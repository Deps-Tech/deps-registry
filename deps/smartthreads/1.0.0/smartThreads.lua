--[[
    smartThreads.lua - Robust Thread Management Library for Lua in MoonLoader/SAMP
    
    Author: https://github.com/DepsCian
    License: MIT
    Version: 1.0.0
    
    DESCRIPTION:
    Solves the critical limitation of Lua in MoonLoader/SAMP where creating more than 3
    concurrent threads causes game crashes. This library implements a thread pool
    with prioritized scheduling to safely manage any number of concurrent tasks.

    EXAMPLES:
    
    Basic usage:
    ```lua
    local st = require('smartThreads')
    st.init()  -- Use default settings
    
    -- Create a simple task
    st.createTask(function()
        wait(1000)  -- Heavy operation
        sampAddChatMessage("Task completed", -1)
    end, "normal", "Notification task")
    ```
    
    Advanced usage:
    ```lua
    local st = require('smartThreads')
    
    -- Custom initialization
    st.init({
        threadLimit = 3,             -- Max concurrent threads
        logLevel = st.LOG_LEVEL.INFO,  -- Only show important logs
        logToFile = true,            -- Write to sampfuncs.log
        logToConsole = false         -- Don't spam chat
    })
    
    -- Delayed task (executes after 5 seconds)
    st.createDelayedTask(function()
        -- Delayed operation
    end, 5.0, "normal", "Delayed task")
    
    -- Periodic task (executes every 30 seconds)
    local timerId = st.createPeriodicTask(function()
        -- Periodic check
    end, 30.0, "low", "Periodic check")
    
    -- Later, stop the periodic task
    st.removePeriodicTask(timerId)
    
    -- Get statistics
    local stats = st.getStats()
    print("Active threads:", stats.activeThreads)
    ```
    
    WARNINGS:
    - DO NOT modify thread limit above 3 unless you absolutely know what you're doing
    - Long-running tasks will block other tasks with same priority, use priorities wisely
    - For tasks requiring UI interaction, use lower priorities to avoid blocking input
    - When working with heavy operations, consider breaking them into smaller tasks
    - Always handle errors inside your task functions - unhandled errors are logged but
      can't be recovered without proper error handling
    - Be cautious with thread-shared resources, as tasks can be executed in any order
]]

local smartThreads = {}
local logger = {}

-- Core configuration
local DEFAULT_THREAD_LIMIT = 3    
local DEFAULT_MANAGER_INTERVAL = 100  
local DEFAULT_TIMEOUT = 10        
local LOG_LEVEL = {               
    DEBUG = 1,
    INFO = 2,
    WARN = 3,
    ERROR = 4,
    NONE = 5
}

-- Thread pool state
local thread_pool = {
    limit = DEFAULT_THREAD_LIMIT,
    active = 0,                   
    tasks = { high = {}, normal = {}, low = {} },
    scheduled = {},              
    periodic = {},               
    active_info = {},            
    stats = { completed = 0, failed = 0 },
    next_id = 0,                 
    manager_active = false       
}

-- Logging configuration
local log_config = {
    level = LOG_LEVEL.INFO,
    to_file = true,
    to_console = false
}

-----------------------------------------------------------------------------
-- Logging system
-----------------------------------------------------------------------------

function logger.init()
    return true
end

function logger.log(level, message, ...)
    if level < log_config.level then return end
    
    local msg = message
    if ... then
        msg = string.format(message, ...)
    end
    
    local level_str = "INFO"
    if level == LOG_LEVEL.DEBUG then level_str = "DEBUG"
    elseif level == LOG_LEVEL.WARN then level_str = "WARN"
    elseif level == LOG_LEVEL.ERROR then level_str = "ERROR"
    end
    
    local log_msg = string.format("[%s] %s", level_str, msg)
    
    if log_config.to_file and type(print) == "function" then
        print(log_msg)
    end
    
    if log_config.to_console and type(sampAddChatMessage) == "function" then
        sampAddChatMessage(log_msg, 0x1E90FF)
    end
    
    return true
end

function logger.debug(message, ...) return logger.log(LOG_LEVEL.DEBUG, message, ...) end
function logger.info(message, ...) return logger.log(LOG_LEVEL.INFO, message, ...) end
function logger.warn(message, ...) return logger.log(LOG_LEVEL.WARN, message, ...) end
function logger.error(message, ...) return logger.log(LOG_LEVEL.ERROR, message, ...) end

-----------------------------------------------------------------------------
-- Utility functions
-----------------------------------------------------------------------------

local function get_unique_id()
    thread_pool.next_id = thread_pool.next_id + 1
    return thread_pool.next_id
end

local function get_queue_size()
    return #thread_pool.tasks.high + #thread_pool.tasks.normal + #thread_pool.tasks.low
end

local function get_next_task()
    if #thread_pool.tasks.high > 0 then
        return table.remove(thread_pool.tasks.high, 1), "high"
    elseif #thread_pool.tasks.normal > 0 then
        return table.remove(thread_pool.tasks.normal, 1), "normal"
    elseif #thread_pool.tasks.low > 0 then
        return table.remove(thread_pool.tasks.low, 1), "low"
    end
    return nil, nil
end

local function validate_priority(priority)
    if priority ~= "high" and priority ~= "normal" and priority ~= "low" then
        return "normal"
    end
    return priority
end

-----------------------------------------------------------------------------
-- Thread management core
-----------------------------------------------------------------------------

-- Creates the background thread that manages task scheduling
local function create_thread_manager()
    if thread_pool.manager_active then 
        return false 
    end
    
    thread_pool.manager_active = true
    
    lua_thread.create(function()
        logger.debug("Thread manager started")
        
        while thread_pool.manager_active do
            local current_time = os.clock()
            local i = 1
            
            -- Process scheduled tasks
            while i <= #thread_pool.scheduled do
                local task = thread_pool.scheduled[i]
                if current_time >= task.execute_at then
                    logger.debug("Scheduled task %s is due for execution", task.id)
                    table.remove(thread_pool.scheduled, i)
                    smartThreads.createTask(task.func, task.priority, task.description)
                else
                    i = i + 1
                end
            end
            
            -- Process periodic tasks
            i = 1
            while i <= #thread_pool.periodic do
                local task = thread_pool.periodic[i]
                if current_time >= task.next_execution then
                    logger.debug("Executing periodic task %s", task.id)
                    task.next_execution = current_time + task.interval
                    smartThreads.createTask(task.func, task.priority, task.description)
                end
                i = i + 1
            end
            
            -- Start tasks from queue if we have capacity
            if thread_pool.active < thread_pool.limit then
                local task, priority = get_next_task()
                if task then
                    execute_task(task, priority)
                end
            end
            
            wait(DEFAULT_MANAGER_INTERVAL)
        end
    end)
    
    return true
end

-- Core function that actually executes a task in a new thread
function execute_task(task, priority)
    if thread_pool.active >= thread_pool.limit then
        -- Prevent task starvation by increasing priority when returning to queue
        local elevated_priority = priority == "low" and "normal" or "high"
        table.insert(thread_pool.tasks[elevated_priority], 1, task)
        return false
    end
    
    thread_pool.active = thread_pool.active + 1
    local task_id = get_unique_id()
    local start_time = os.clock()
    
    thread_pool.active_info[task_id] = {
        id = task_id,
        description = task.description or "Unnamed Task",
        priority = priority,
        start_time = start_time
    }
    
    logger.debug("Starting task #%d (%s) [%s]", 
                task_id, thread_pool.active_info[task_id].description, priority)
    
    lua_thread.create(function()
        -- Critical: prevents the "cannot resume non-suspended coroutine" error
        wait(0)
        
        local success, err = pcall(function()
            task.func()
        end)
        
        local execution_time = os.clock() - start_time
        
        if success then
            thread_pool.stats.completed = thread_pool.stats.completed + 1
            logger.debug("Task #%d completed in %.2f seconds", task_id, execution_time)
        else
            thread_pool.stats.failed = thread_pool.stats.failed + 1
            logger.error("Task #%d failed: %s", task_id, tostring(err))
        end
        
        thread_pool.active_info[task_id] = nil
        thread_pool.active = thread_pool.active - 1
    end)
    
    return true, task_id
end

-----------------------------------------------------------------------------
-- Public API
-----------------------------------------------------------------------------

function smartThreads.init(options)
    options = options or {}
    
    thread_pool.limit = options.threadLimit or DEFAULT_THREAD_LIMIT
    log_config.level = options.logLevel or LOG_LEVEL.INFO
    log_config.to_file = options.logToFile ~= false
    log_config.to_console = options.logToConsole or false
    
    logger.init()
    create_thread_manager()
    
    logger.info("smartThreads initialized with thread limit: %d", thread_pool.limit)
    return true
end

function smartThreads.createTask(func, priority, description)
    if type(func) ~= "function" then
        logger.error("createTask: 'func' must be a function")
        return false
    end
    
    priority = validate_priority(priority or "normal")
    
    local task = {
        func = func,
        description = description or "Task " .. get_unique_id()
    }
    
    if thread_pool.active < thread_pool.limit then
        return execute_task(task, priority)
    else
        table.insert(thread_pool.tasks[priority], task)
        logger.debug("Task queued with priority [%s]", priority)
        return false, get_queue_size()
    end
end

function smartThreads.createDelayedTask(func, delay, priority, description)
    if type(func) ~= "function" then
        logger.error("createDelayedTask: 'func' must be a function")
        return false
    end
    
    delay = tonumber(delay) or 1
    priority = validate_priority(priority or "normal")
    
    local task_id = get_unique_id()
    local task = {
        id = task_id,
        func = func,
        execute_at = os.clock() + delay,
        priority = priority,
        description = description or "Delayed Task " .. task_id
    }
    
    table.insert(thread_pool.scheduled, task)
    logger.debug("Delayed task #%d created (delay: %.1f sec)", task_id, delay)
    
    return task_id
end

function smartThreads.createPeriodicTask(func, interval, priority, description)
    if type(func) ~= "function" then
        logger.error("createPeriodicTask: 'func' must be a function")
        return false
    end
    
    interval = tonumber(interval) or 1
    priority = validate_priority(priority or "normal")
    
    local task_id = get_unique_id()
    local task = {
        id = task_id,
        func = func,
        interval = interval,
        next_execution = os.clock() + interval,
        priority = priority,
        description = description or "Periodic Task " .. task_id
    }
    
    table.insert(thread_pool.periodic, task)
    logger.debug("Periodic task #%d created (interval: %.1f sec)", task_id, interval)
    
    return task_id
end

function smartThreads.removePeriodicTask(id)
    for i, task in ipairs(thread_pool.periodic) do
        if task.id == id then
            table.remove(thread_pool.periodic, i)
            logger.debug("Periodic task #%d removed", id)
            return true
        end
    end
    logger.warn("Periodic task #%d not found for removal", id)
    return false
end

-- Waits until all queued tasks are completed
-- WARNING: This will block the calling thread, use with caution
function smartThreads.waitForQueue(timeout)
    local start_time = os.clock()
    timeout = timeout or DEFAULT_TIMEOUT
    
    logger.debug("Waiting for queue to empty (timeout: %d sec)", timeout)
    
    while get_queue_size() > 0 or thread_pool.active > 0 do
        if timeout and os.clock() - start_time > timeout then
            logger.warn("Queue wait timed out after %d seconds", timeout)
            return false
        end
        wait(100)
    end
    
    logger.debug("All queued tasks completed")
    return true
end

function smartThreads.clearQueues()
    thread_pool.tasks.high = {}
    thread_pool.tasks.normal = {}
    thread_pool.tasks.low = {}
    thread_pool.scheduled = {}
    
    logger.info("All task queues cleared")
    return true
end

function smartThreads.getStats()
    return {
        activeThreads = thread_pool.active,
        maxThreads = thread_pool.limit,
        queueSizes = {
            high = #thread_pool.tasks.high,
            normal = #thread_pool.tasks.normal,
            low = #thread_pool.tasks.low,
            total = get_queue_size()
        },
        scheduledTasks = #thread_pool.scheduled,
        periodicTasks = #thread_pool.periodic,
        completed = thread_pool.stats.completed,
        failed = thread_pool.stats.failed
    }
end

function smartThreads.getActiveTasks()
    local result = {}
    for id, info in pairs(thread_pool.active_info) do
        table.insert(result, {
            id = id,
            description = info.description,
            priority = info.priority,
            runtime = os.clock() - info.start_time
        })
    end
    return result
end

-- CAUTION: Modifying thread limit above 3 may cause game crashes!
function smartThreads.setThreadLimit(limit)
    if type(limit) ~= "number" or limit < 1 then
        logger.error("setThreadLimit: 'limit' must be a positive number")
        return false
    end
    
    -- Warning if setting above the safe threshold
    if limit > 3 then
        logger.warn("!!!WARNING!!! Setting thread limit above 3 (%d) may cause game crashes!", limit)
    end
    
    thread_pool.limit = limit
    logger.info("Thread limit changed to %d", limit)
    return true
end

function smartThreads.setLogLevel(level)
    if type(level) ~= "number" or level < LOG_LEVEL.DEBUG or level > LOG_LEVEL.NONE then
        logger.error("setLogLevel: Invalid log level")
        return false
    end
    
    log_config.level = level
    logger.info("Log level set to %d", level)
    return true
end

function smartThreads.configureLogging(options)
    if options.level then
        log_config.level = options.level
    end
    if options.toFile ~= nil then
        log_config.to_file = options.toFile
    end
    if options.toConsole ~= nil then
        log_config.to_console = options.toConsole
    end
    
    logger.info("Logging configuration updated")
    return true
end

function smartThreads.printStats()
    local stats = smartThreads.getStats()
    
    logger.info("=== smartThreads Statistics ===")
    logger.info("Active threads: %d/%d", stats.activeThreads, stats.maxThreads)
    logger.info("Queue sizes: %d total (High: %d | Normal: %d | Low: %d)", 
                stats.queueSizes.total, stats.queueSizes.high, 
                stats.queueSizes.normal, stats.queueSizes.low)
    logger.info("Scheduled tasks: %d | Periodic tasks: %d", 
                stats.scheduledTasks, stats.periodicTasks)
    logger.info("Completed tasks: %d | Failed tasks: %d", 
                stats.completed, stats.failed)
    
    local active_tasks = smartThreads.getActiveTasks()
    if #active_tasks > 0 then
        logger.info("Active tasks:")
        for i, task in ipairs(active_tasks) do
            logger.info("  #%d: %s [%s] (%.1f sec)", 
                      task.id, task.description, 
                      task.priority, task.runtime)
        end
    end
    
    return stats
end

function smartThreads.shutdown()
    logger.info("Shutting down smartThreads")
    thread_pool.manager_active = false
    smartThreads.clearQueues()
    logger.info("smartThreads shutdown complete")
    return true
end

smartThreads.LOG_LEVEL = LOG_LEVEL

return smartThreads 