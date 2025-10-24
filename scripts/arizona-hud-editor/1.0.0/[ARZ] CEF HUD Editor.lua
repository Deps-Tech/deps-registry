script_author('Willy4ka')
script_name('Arizona HUD Editor')
script_url('https://www.blast.hk/threads/192708/')

local json = require('carbJsonConfig')
local ffi = require('ffi')
local imgui = require('mimgui')
local encoding = require('encoding')
local lfs = require('lfs')
encoding.default = 'CP1251'
local u8 = encoding.UTF8

local inicfg = require 'inicfg'
local directIni = 'arzhud.ini'
local ini = inicfg.load(inicfg.load({
    main = {
        config_selected = 'default.json'
    },
}, directIni))
inicfg.save(ini, directIni)

local d = 'default.json'
local script_path = getWorkingDirectory()..'\\CEFHUD\\'
if not doesDirectoryExist(script_path) then createDirectory(script_path) end

local renderWindow = imgui.new.bool(false)
local createConfig = imgui.new.char[64]()
local last_sended = os.clock()

local settings = {
    enabled = imgui.new.bool(true),
    project = imgui.new.char[64](),
    title = imgui.new.char[64](),
    type = imgui.new.char[64](),
    flag = imgui.new.int(0),
    id = imgui.new.int(0),
    logo = imgui.new.int(0),
    multiplier = imgui.new.int(0),
    flagurl = imgui.new.char[512](),
    logourl = imgui.new.char[512](),
    background = imgui.new.char[512](),
    backgroundRouding = imgui.new.int(),
    weaponBackgroundRouding = imgui.new.int(),
    search = imgui.new.char[64](),
    deleteOnline = imgui.new.bool(false),
    deleteId = imgui.new.bool(false),
    deleteServerId = imgui.new.bool(false),
    deletePromotion = imgui.new.bool(false),
    deleteBackground = imgui.new.bool(false),
    deleteWeaponBackground = imgui.new.bool(false),
    deleteADDVIP = imgui.new.bool(false),
    divShadow = imgui.new.bool(false),
    divShadowX = imgui.new.int(0),
    divShadowY = imgui.new.int(0),
    divShadowBlur = imgui.new.int(0),
    divShadowSize = imgui.new.int(0),
    textShadow = imgui.new.bool(false),
    textShadowX = imgui.new.int(0),
    textShadowY = imgui.new.int(0),
    textShadowBlur = imgui.new.int(0),
    moneyShadow = imgui.new.bool(false),
    moneyShadowX = imgui.new.int(0),
    moneyShadowY = imgui.new.int(0),
    moneyShadowBlur = imgui.new.int(0),
    moneyShadowSize = imgui.new.int(0),
    weapons = {},
    colors = {
        health = imgui.new.float[4](1.00, 0.11, 0.22, 1.00),
        armour = imgui.new.float[4](1.00, 1.00, 1.00, 1.00),
        satiety = imgui.new.float[4](0.95, 0.54, 0.15, 1.00),
        money = imgui.new.float[4](1.00, 1.00, 1.00, 1.00),
        moneyIcon = imgui.new.float[4](0.40, 1.00, 0.39, 1.00),
        background = imgui.new.float[4](0.01, 0.02, 0.06, 0.85),
        divShadow = imgui.new.float[4](0.00, 0.00, 0.00, 1.00),
        textShadow = imgui.new.float[4](0.00, 0.00, 0.00, 1.00),
        moneyShadow = imgui.new.float[4](0.00, 0.00, 0.00, 1.00),
        serverNumber = imgui.new.float[4](1.00, 1.00, 1.00, 0.28),
        plusMoney = imgui.new.float[4](0.46, 0.82, 0.33, 1.00),
        minusMoney = imgui.new.float[4](1.00, 0.13, 0.13, 1.00),
        addVipGr1 = imgui.new.float[4](0.95, 0.75, 0.08, 1.00),
        addVipGr2 = imgui.new.float[4](1.00, 0.87, 0.48, 1.00),
        addVipText = imgui.new.float[4](0.27, 0.23, 0.10, 1.00),
        greenZoneText = imgui.new.float[4](0.46, 0.82, 0.33, 1.00),
        greenZoneText2 = imgui.new.float[4](1.00, 1.00, 1.00, 1.00),
        greenZoneBg = imgui.new.float[4](0.01, 0.02, 0.06, 1.00),
        greenZoneGr1 = imgui.new.float[4](0.38, 0.87, 0.33, 1.00),
        greenZoneGr2 = imgui.new.float[4](0.38, 0.87, 0.33, 0.00),
        greenZoneIcon = imgui.new.float[4](0.38, 0.87, 0.33, 1.00),
        X4 = imgui.new.float[4](1.00, 0.76, 0.17, 1.00),
    },
}
function getValidWeapons()
    local result = {}
    local FLAWeaponConfigPath = getGameDirectory()..'\\data\\gtasa_weapon_config.dat'
    if doesFileExist(FLAWeaponConfigPath) then
        for line in io.lines(FLAWeaponConfigPath) do
            if line:find('^%s*(%d+)%s*([A-Za-z0-9_]+)%s*') then
                local id, name = line:match('^%s*(%d+)%s*([A-Za-z0-9_]+)%s*')
                table.insert(result, {id = tonumber(id), name = name})
            end
        end
    end
    return result
end
local allweapons = getValidWeapons()

for k, v in pairs(allweapons) do
    settings.weapons[v.name] = imgui.new.char[512]()
end

local config_path = script_path..ini.main.config_selected
json.load(config_path, settings)

local hud_data = {
    project = '',
    title = '',
    type = '',
    flag = 0,
    id = 0,
    logo = 0,
    multiplier = 0,
}

local old_hud = hud_data

imgui.OnInitialize(function()
    SoftBlueTheme()
    imgui.GetIO().IniFilename = nil
end)

local newFrame = imgui.OnFrame(
    function() return renderWindow[0] end,
    function()
        local resX, resY = getScreenResolution()
        local sizeX, sizeY = 1000, 650
        imgui.SetNextWindowPos(imgui.ImVec2(resX / 2, resY / 2), imgui.Cond.FirstUseEver, imgui.ImVec2(0.5, 0.5))
        imgui.SetNextWindowSize(imgui.ImVec2(sizeX, sizeY), imgui.Cond.FirstUseEver)
        if imgui.Begin('Arizona HUD Editor by Willy4ka', renderWindow, imgui.WindowFlags.NoCollapse + imgui.WindowFlags.NoResize + imgui.WindowFlags.NoScrollbar) then
            imgui.BeginChild('##1', imgui.ImVec2(490, -1), true)
            imgui.Text(u8('После перезагрузки скрипта откройте и закройте меню доната'))
            if imgui.Checkbox(u8('Включить'), settings.enabled) then save() end
            imgui.TextWrapped(u8('В строке URL указывайте только название файла, сами файлы должны храниться в папке moonloader\\CEFHUD'))
            imgui.TextColored(imgui.ImVec4(0.70, 0.70, 0.90, 1.00), u8('Ссылка на тему'))
            if imgui.IsItemClicked(0) then
                os.execute("explorer "..script.this.url)
            end
            if imgui.CollapsingHeader(u8("Редактирование информации о сервере")) then
                if imgui.InputTextWithHint('##server', u8('Название сервера'), settings.title, ffi.sizeof(settings.title)) then save() end
                if imgui.InputTextWithHint('##project', u8('Название проекта'), settings.project, ffi.sizeof(settings.project)) then save() end
                if imgui.InputTextWithHint('##type', u8('Тип проекта'), settings.type, ffi.sizeof(settings.type)) then save() end
                imgui.NewLine()
                if imgui.InputInt(u8('Номер сервера'), settings.id) then save() end
                if imgui.InputInt(u8('Флаг сервера'), settings.flag) then save() end
                if imgui.InputInt(u8('Логотип проекта'), settings.logo) then save() end
                if imgui.InputInt(u8('Множитель сервера'), settings.multiplier) then save() end
            end

            if imgui.CollapsingHeader(u8("Удаление элементов")) then
                if imgui.Checkbox(u8('Удалять номер сервера'), settings.deleteServerId) then save() end
                if imgui.Checkbox(u8('Удалять онлайн сервера'), settings.deleteOnline) then save() end
                if imgui.Checkbox(u8('Удалять свой id'), settings.deleteId) then save() end
                if imgui.Checkbox(u8('Удалять акции (X4)'), settings.deletePromotion) then save() end
                if imgui.Checkbox(u8('Удалять фон'), settings.deleteBackground) then save() end
                if imgui.Checkbox(u8('Удалять фон оружия'), settings.deleteWeaponBackground) then save() end
                if imgui.Checkbox(u8('Удалять ADD VIP'), settings.deleteADDVIP) then save() end
            end

            if imgui.CollapsingHeader(u8("Редактирование изображений")) then
                imgui.Text(u8('URL картинки для замены логотипа и флага'))
                if imgui.InputTextWithHint(u8('Логотип'), u8('Логотип'), settings.logourl, ffi.sizeof(settings.logourl)) then save() end
                if imgui.InputTextWithHint(u8('Флаг'), u8('Флаг'), settings.flagurl, ffi.sizeof(settings.flagurl)) then save() end
                imgui.NewLine()
                if imgui.InputTextWithHint(u8('Фон'), u8('Фон'), settings.background, ffi.sizeof(settings.background)) then save() end
                if imgui.InputInt(u8('Скругление фона'), settings.backgroundRouding) then save() end
                if imgui.InputInt(u8('Скругление фона\nоружия'), settings.weaponBackgroundRouding) then save() end
            end

            if imgui.CollapsingHeader(u8("Редактирование цветов элементов")) then
                if imgui.ColorEdit4(u8('Фон'), settings.colors.background, imgui.ColorEditFlags.AlphaBar) then save() end
                if imgui.ColorEdit4(u8('Здоровье'), settings.colors.health, imgui.ColorEditFlags.AlphaBar) then save() end
                if imgui.ColorEdit4(u8('Броня'), settings.colors.armour, imgui.ColorEditFlags.AlphaBar) then save() end
                if imgui.ColorEdit4(u8('Сытость'), settings.colors.satiety, imgui.ColorEditFlags.AlphaBar) then save() end
                if imgui.ColorEdit4(u8('Деньги'), settings.colors.money, imgui.ColorEditFlags.AlphaBar) then save() end
                if imgui.ColorEdit4(u8('Прибавление денег'), settings.colors.plusMoney, imgui.ColorEditFlags.AlphaBar) then save() end
                if imgui.ColorEdit4(u8('Убавление денег'), settings.colors.minusMoney, imgui.ColorEditFlags.AlphaBar) then save() end
                if imgui.ColorEdit4(u8('Иконка денег'), settings.colors.moneyIcon, imgui.ColorEditFlags.AlphaBar) then save() end
                if imgui.ColorEdit4(u8('Номер сервера'), settings.colors.serverNumber, imgui.ColorEditFlags.AlphaBar) then save() end
                if imgui.ColorEdit4(u8('Градиент ADD VIP 1'), settings.colors.addVipGr1, imgui.ColorEditFlags.AlphaBar) then save() end
                if imgui.ColorEdit4(u8('Градиент ADD VIP 2'), settings.colors.addVipGr2, imgui.ColorEditFlags.AlphaBar) then save() end
                if imgui.ColorEdit4(u8('Текст ADD VIP'), settings.colors.addVipText, imgui.ColorEditFlags.AlphaBar) then save() end
                if imgui.ColorEdit4(u8('Фон зз'), settings.colors.greenZoneBg, imgui.ColorEditFlags.AlphaBar) then save() end
                if imgui.ColorEdit4(u8('Градиент зз 1'), settings.colors.greenZoneGr1, imgui.ColorEditFlags.AlphaBar) then save() end
                if imgui.ColorEdit4(u8('Градиент зз 2'), settings.colors.greenZoneGr2, imgui.ColorEditFlags.AlphaBar) then save() end
                if imgui.ColorEdit4(u8('Иконка зз'), settings.colors.greenZoneIcon, imgui.ColorEditFlags.AlphaBar) then save() end
                if imgui.ColorEdit4(u8('Текст зз 1'), settings.colors.greenZoneText, imgui.ColorEditFlags.AlphaBar) then save() end
                if imgui.ColorEdit4(u8('Текст зз 2'), settings.colors.greenZoneText2, imgui.ColorEditFlags.AlphaBar) then save() end
                if imgui.ColorEdit4(u8('Акции (X4)'), settings.colors.X4, imgui.ColorEditFlags.AlphaBar) then save() end
            end

            if imgui.CollapsingHeader(u8("Редактирование теней элементов")) then
                if imgui.Checkbox(u8('Редактирование тени индикаторов'), settings.divShadow) then save() end
                if imgui.ColorEdit4(u8('Цвет тени##1'), settings.colors.divShadow, imgui.ColorEditFlags.AlphaBar) then save() end
                if imgui.InputInt(u8('Размер тени##1'), settings.divShadowSize) then save() end
                if imgui.InputInt(u8('Блюр тени##1'), settings.divShadowBlur) then save() end
                if imgui.InputInt(u8('Смещение тени X##1'), settings.divShadowX) then save() end
                if imgui.InputInt(u8('Смещение тени Y##1'), settings.divShadowY) then save() end

                if imgui.Checkbox(u8('Редактирование тени иконки денег'), settings.moneyShadow) then save() end
                if imgui.ColorEdit4(u8('Цвет тени##2'), settings.colors.moneyShadow, imgui.ColorEditFlags.AlphaBar) then save() end
                if imgui.InputInt(u8('Размер тени##2'), settings.moneyShadowSize) then save() end
                if imgui.InputInt(u8('Блюр тени##2'), settings.moneyShadowBlur) then save() end
                if imgui.InputInt(u8('Смещение тени X##2'), settings.moneyShadowX) then save() end
                if imgui.InputInt(u8('Смещение тени Y##2'), settings.moneyShadowY) then save() end

                if imgui.Checkbox(u8('Редактирование тени кол-ва денег'), settings.textShadow) then save() end
                if imgui.ColorEdit4(u8('Цвет тени##3'), settings.colors.textShadow, imgui.ColorEditFlags.AlphaBar) then save() end
                if imgui.InputInt(u8('Блюр тени##3'), settings.textShadowBlur) then save() end
                if imgui.InputInt(u8('Смещение тени X##3'), settings.textShadowX) then save() end
                if imgui.InputInt(u8('Смещение тени Y##3'), settings.textShadowY) then save() end
            end
            if imgui.CollapsingHeader(u8('Загрузка конфига')) then
                imgui.BeginChild('##config_default', imgui.ImVec2(-1, 60), true)
                imgui.SetCursorPosY(imgui.GetCursorPosY()+7)
                imgui.Text(u8(d))
                imgui.SameLine()
                imgui.SetCursorPosY(imgui.GetCursorPosY()-7)
                if imgui.Button(u8('Загрузить')) then
                    ini.main.config_selected = d
                    inicfg.save(ini, directIni)
                    config_path = script_path..ini.main.config_selected
                    json.load(config_path, settings)
                    settings("reset")
                    save()
                end
                imgui.EndChild()
                for k, v in pairs(getAllConfigs()) do
                    imgui.BeginChild('##config_'..v, imgui.ImVec2(-1, 60), true)
                    imgui.SetCursorPosY(imgui.GetCursorPosY()+7)
                    imgui.Text(u8(v))
                    imgui.SameLine()
                    imgui.SetCursorPosY(imgui.GetCursorPosY()-7)
                    if imgui.Button(u8('Загрузить')) then
                        ini.main.config_selected = v
                        inicfg.save(ini, directIni)
                        config_path = script_path..ini.main.config_selected
                        json.load(config_path, settings)
                        save()
                    end
                    imgui.SameLine()
                    if imgui.Button(u8('Сохранить')) then
                        settings()
                    end
                    imgui.SameLine()
                    if imgui.Button(u8('Удалить')) then
                        os.remove(script_path..v)
                        if config_path == script_path..v then
                            ini.main.config_selected = d
                            config_path = script_path..d
                            json.load(config_path, settings)
                            inicfg.save(ini, directIni)
                        end
                    end
                    imgui.EndChild()
                end
                imgui.BeginChild('##create_config', imgui.ImVec2(-1, 60), true)
                imgui.InputTextWithHint('##create_config_input', u8('Название конфига'), createConfig, ffi.sizeof(createConfig))
                imgui.SameLine()
                if imgui.Button(u8('Создать')) then
                    for k, v in pairs(getAllConfigs()) do
                        if v == u8:decode(ffi.string(createConfig))..'.json' then return end
                    end
                    config_path = script_path..u8:decode(ffi.string(createConfig))..'.json'
                    json.save(config_path, settings)
                end
                imgui.EndChild()
            end

            imgui.EndChild()
            imgui.SameLine()
            imgui.BeginChild('##2', imgui.ImVec2(-1, -1), true)
            imgui.Text(u8('URL картинки для замены оружия'))

            imgui.InputTextWithHint(u8('Поиск'), u8('Поиск'), settings.search, ffi.sizeof(settings.search))
            imgui.NewLine()
            for k, v in pairs(allweapons) do
                if #ffi.string(settings.search) > 0 then
                    if u8(v.name):lower():find(ffi.string(settings.search):lower()) then
                        if imgui.InputTextWithHint(u8(v.name)..'##'..v.id, u8(v.name), settings.weapons[v.name], ffi.sizeof(settings.weapons[v.name])) then save() end
                    end
                else
                    if imgui.InputTextWithHint(u8(v.name)..'##'..v.id, u8(v.name), settings.weapons[v.name], ffi.sizeof(settings.weapons[v.name])) then save() end
                end
            end
            imgui.EndChild()
            imgui.End()
        end
    end
)

function main()
    while not isSampAvailable() do wait(0) end
    sampRegisterChatCommand('chud', function()
        renderWindow[0] = not renderWindow[0]
    end)
    while true do
        wait(0)
        if settings.enabled[0] then
            setImages()
        end
    end
end
function onReceivePacket(id, bs)
    if id == 220 then
        raknetBitStreamReadInt8(bs);
        local ptype = raknetBitStreamReadInt8(bs)
        if ptype == 17 then
            raknetBitStreamReadInt32(bs)
            local length = raknetBitStreamReadInt16(bs)
            local encoded = raknetBitStreamReadInt8(bs)
            if length > 0 then
                local text = (encoded ~= 0) and raknetBitStreamDecodeString(bs, length + encoded) or raknetBitStreamReadString(bs, length)
                local event, data = text:match('window%.executeEvent%(\'(.+)\',%s*`%[(.+)%]`%);');
                if event == 'event.arizonahud.serverInfo' then
                    hud_data, old_hud = decodeJson(data), decodeJson(data)
                    save()
                elseif event == 'event.arizonahud.setRadialKey' then
                    save()
                end
            end
        end
    end
end

function set(v, u, p, i)
    if #v > 0 then
        if doesDirectoryExist(getWorkingDirectory()..'\\CEFHUD') then
            if not v:find('https://') then
                v = 'file:///'..getWorkingDirectory()..'/CEFHUD/'..v
                v = v:gsub('\\', '/'):gsub('%s', '%%20')
            end
        end
        if u then
            v = ('url(%s)'):format(v)
        end
        if i then
            evalanon(([[
                if(document.querySelector('%s') != null)
                    if(document.querySelector('%s').src != '%s')
                        document.querySelector('%s').src = '%s'
            ]]):format(p, p, v, p, v))
        else
            evalanon(([[
                if(document.querySelector('%s') != null)
                    if(document.querySelector('%s').style.backgroundImage != '%s')
                        document.querySelector('%s').style.backgroundImage = '%s'
            ]]):format(p, p, v, p, v))
        end
    end
end

addEventHandler('onSendPacket', function (id, bs, priority, reliability, orderingChannel)
    if id == 220 then
        raknetBitStreamReadInt8(bs);
        if raknetBitStreamReadInt8(bs) == 17 then
            if raknetBitStreamReadInt32(bs) == 0 then
                if os.clock() - last_sended < 1.5 or sampGetGamestate() ~= 3 then
                    return false
                end
                last_sended = os.clock()
            end
        end
    end
end)

function evalanon(code)
    evalcef(("(() => {%s})()"):format(code))
end

function evalcef(code, encoded)
    encoded = encoded or 0
    local bs = raknetNewBitStream();
    raknetBitStreamWriteInt8(bs, 17);
    raknetBitStreamWriteInt32(bs, 0);
    raknetBitStreamWriteInt16(bs, #code);
    raknetBitStreamWriteInt8(bs, encoded);
    raknetBitStreamWriteString(bs, code);
    raknetEmulPacketReceiveBitStream(220, bs);
    raknetDeleteBitStream(bs);
end

function save()
    if settings.enabled[0] and hud_data then
        hud_data.project = u8:decode(check(ffi.string(settings.project)) or old_hud.project)
        hud_data.title = u8:decode(check(ffi.string(settings.title)) or old_hud.title)
        hud_data.type = u8:decode(check(ffi.string(settings.type)) or old_hud.type)
        hud_data.flag = check(settings.flag[0]) or old_hud.flag
        hud_data.id = check(settings.id[0]) or old_hud.id
        hud_data.logo = check(settings.logo[0]) or old_hud.logo
        hud_data.multiplier = check(settings.multiplier[0]) or old_hud.multiplier

        local colorHealth = string.format('rgba(%d, %d, %d, %f)', settings.colors.health[0]*255, settings.colors.health[1]*255, settings.colors.health[2]*255, settings.colors.health[3])
        local colorArmour = string.format('rgba(%d, %d, %d, %f)', settings.colors.armour[0]*255, settings.colors.armour[1]*255, settings.colors.armour[2]*255, settings.colors.armour[3])
        local colorSatiety = string.format('rgba(%d, %d, %d, %f)', settings.colors.satiety[0]*255, settings.colors.satiety[1]*255, settings.colors.satiety[2]*255, settings.colors.satiety[3])
        local colorMoney = string.format('rgba(%d, %d, %d, %f)', settings.colors.money[0]*255, settings.colors.money[1]*255, settings.colors.money[2]*255, settings.colors.money[3])
        local colorMoneyIcon = string.format('rgba(%d, %d, %d, %f)', settings.colors.moneyIcon[0]*255, settings.colors.moneyIcon[1]*255, settings.colors.moneyIcon[2]*255, settings.colors.moneyIcon[3])
        local colorBackground = string.format('rgba(%d, %d, %d, %f)', settings.colors.background[0]*255, settings.colors.background[1]*255, settings.colors.background[2]*255, settings.colors.background[3])
        local colordivShadow = string.format('rgba(%d, %d, %d, %f)', settings.colors.divShadow[0]*255, settings.colors.divShadow[1]*255, settings.colors.divShadow[2]*255, settings.colors.divShadow[3])
        local colortextShadow = string.format('rgba(%d, %d, %d, %f)', settings.colors.textShadow[0]*255, settings.colors.textShadow[1]*255, settings.colors.textShadow[2]*255, settings.colors.textShadow[3])
        local colormoneyShadow = string.format('rgba(%d, %d, %d, %f)', settings.colors.moneyShadow[0]*255, settings.colors.moneyShadow[1]*255, settings.colors.moneyShadow[2]*255, settings.colors.moneyShadow[3])
        local colorServerNumber = string.format('rgba(%d, %d, %d, %f)', settings.colors.serverNumber[0]*255, settings.colors.serverNumber[1]*255, settings.colors.serverNumber[2]*255, settings.colors.serverNumber[3])
        local colorAddVIP1 = string.format('rgba(%d, %d, %d, %f)', settings.colors.addVipGr1[0]*255, settings.colors.addVipGr1[1]*255, settings.colors.addVipGr1[2]*255, settings.colors.addVipGr1[3])
        local colorAddVIP2 = string.format('rgba(%d, %d, %d, %f)', settings.colors.addVipGr2[0]*255, settings.colors.addVipGr2[1]*255, settings.colors.addVipGr2[2]*255, settings.colors.addVipGr2[3])
        local colorAddVIPText = string.format('rgba(%d, %d, %d, %f)', settings.colors.addVipText[0]*255, settings.colors.addVipText[1]*255, settings.colors.addVipText[2]*255, settings.colors.addVipText[3])
        local colorX4 = string.format('rgba(%d, %d, %d, %f)', settings.colors.X4[0]*255, settings.colors.X4[1]*255, settings.colors.X4[2]*255, settings.colors.X4[3])
        evalanon(([[
            let addvip = `<svg class="player-info__vip-logo" width="88" height="16" viewBox="0 0 88 16" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M0.48 15L6.48 0.899999H9.32L15.32 15H12.1L10.82 11.86H4.9L3.62 15H0.48ZM6 9.14H9.72L7.86 4.6L6 9.14ZM17.4808 15V0.999999H22.9408C25.1274 0.999999 26.9141 1.66667 28.3008 3C29.6874 4.32 30.3808 5.98667 30.3808 8C30.3808 10 29.6808 11.6667 28.2808 13C26.8941 14.3333 25.1141 15 22.9408 15H17.4808ZM20.5608 12.22H22.9408C24.2074 12.22 25.2274 11.8333 26.0008 11.06C26.7741 10.2733 27.1608 9.25333 27.1608 8C27.1608 6.76 26.7674 5.74667 25.9808 4.96C25.2074 4.17333 24.1941 3.78 22.9408 3.78H20.5608V12.22ZM33.1253 15V0.999999H38.5853C40.772 0.999999 42.5586 1.66667 43.9453 3C45.332 4.32 46.0253 5.98667 46.0253 8C46.0253 10 45.3253 11.6667 43.9253 13C42.5386 14.3333 40.7586 15 38.5853 15H33.1253ZM36.2053 12.22H38.5853C39.852 12.22 40.872 11.8333 41.6453 11.06C42.4186 10.2733 42.8053 9.25333 42.8053 8C42.8053 6.76 42.412 5.74667 41.6253 4.96C40.852 4.17333 39.8386 3.78 38.5853 3.78H36.2053V12.22ZM59.2259 15.1L53.5659 0.999999H56.9659L60.6259 10.86L64.2859 0.999999H67.6059L61.9459 15.1H59.2259ZM69.9059 15V0.999999H72.9859V15H69.9059ZM76.4847 15V0.999999H82.2047C83.858 0.999999 85.1647 1.44667 86.1247 2.34C87.0847 3.23333 87.5647 4.41333 87.5647 5.88C87.5647 7.44 87.0314 8.65333 85.9647 9.52C84.898 10.3733 83.5447 10.8 81.9047 10.8H79.5647V15H76.4847ZM79.5647 8.06H82.0047C82.7647 8.06 83.358 7.86 83.7847 7.46C84.2247 7.06 84.4447 6.54 84.4447 5.9C84.4447 5.22 84.2247 4.7 83.7847 4.34C83.3447 3.96667 82.7314 3.78 81.9447 3.78H79.5647V8.06Z" fill="url(#paint0_linear_15689_113)"/>
            <defs>
            <linearGradient id="paint0_linear_15689_113" x1="-6.63102e-07" y1="9.5" x2="89" y2="9.5" gradientUnits="userSpaceOnUse">
            <stop stop-color="%s"/>
            <stop offset="1" stop-color="%s"/>
            </linearGradient>
            </defs>
            </svg>`
            if (document.querySelector(`.player-info__vip-logo`)) document.querySelector(`.player-info__vip-logo`).remove()
            document.querySelector(`.player-info__vip`).insertAdjacentHTML("afterbegin", addvip)
            document.querySelector('.player-info__vip-status').style.background = "linear-gradient(90deg,%s,%s)"
            document.querySelector('.player-info__vip-status-duration').style.color = "%s"
            document.querySelector('.player-info__server-info-bar-multiplyer').style.backgroundImage = "radial-gradient(circle,%s 30%%,transparent 90%%)"
            document.querySelector('.player-info__server-info-bar-multiplyer').style.borderColor = "%s"
        ]]):format(colorAddVIP1, colorAddVIP2, colorAddVIP1, colorAddVIP2, colorAddVIPText, colorX4, colorX4))
        evalanon(([[
            let el = document.querySelector('.arizona-hud__player-info > .player-info > .player-info__params > .player-info__indicators')
            if (el != null)
            {
                el.children[0].querySelector('.circle-indicator').style.setProperty('--icon-color', '%s')
                el.children[1].querySelector('.circle-indicator').style.setProperty('--icon-color', '%s')
                el.children[2].querySelector('.circle-indicator').style.setProperty('--icon-color', '%s')
                document.querySelector('.player-info__total').style.color = '%s'
                document.querySelector('.player-info__dollar-icon').style.background = '%s'
                document.querySelector('.player-info__params').style.background = '%s'
                document.querySelector('.player-info__server-number').style.background = '%s'
            }
        ]]):format(colorHealth, colorArmour, colorSatiety, colorMoney, colorMoneyIcon, colorBackground, colorServerNumber))
        if settings.divShadow[0] then
            evalanon(([[
                let x = '%dpx'
                let y = '%dpx'
                let b = '%dpx'
                let s = '%dpx'
                let c = '%s'
                let el = document.querySelector('.arizona-hud__player-info > .player-info > .player-info__params > .player-info__indicators')
                if (el != null)
                {
                    el.children[0].querySelector('.circle-indicator').style.boxShadow = `${x} ${y} ${b} ${s} ${c}`
                    el.children[1].querySelector('.circle-indicator').style.boxShadow = `${x} ${y} ${b} ${s} ${c}`
                    el.children[2].querySelector('.circle-indicator').style.boxShadow = `${x} ${y} ${b} ${s} ${c}`
                }
            ]]):format(settings.divShadowX[0], settings.divShadowY[0], settings.divShadowBlur[0], settings.divShadowSize[0], colordivShadow))
        end
        if settings.moneyShadow[0] then
            evalanon(([[
                if(document.querySelector('.player-info__dollar-icon') != null)
                    document.querySelector('.player-info__dollar-icon').style.boxShadow = '%dpx %dpx %dpx %dpx %s'
            ]]):format(settings.moneyShadowX[0], settings.moneyShadowY[0], settings.moneyShadowBlur[0], settings.moneyShadowSize[0], colormoneyShadow))
        end
        if settings.textShadow[0] then
            evalanon(([[
                if(document.querySelector('.player-info__total') != null)
                    document.querySelector('.player-info__total').style.textShadow = '%s %dpx %dpx %dpx'
            ]]):format(colortextShadow, settings.textShadowX[0], settings.textShadowY[0], settings.textShadowBlur[0]))
        end

        if settings.deleteServerId[0] then evalanon [[
                if(document.querySelector('.player-info__server-number') != null)
                    document.querySelector('.player-info__server-number').remove()
            ]]
        end
        if settings.deleteOnline[0] then evalanon [[
                if(document.querySelector('.player-info__users-online') != null)
                    document.querySelector('.player-info__users-online').remove()
            ]]
        end
        if settings.deleteId[0] then evalanon [[
                if(document.querySelector('.player-info__user-id') != null)
                    document.querySelector('.player-info__user-id').remove()
            ]]
        end
        if settings.deletePromotion[0] then evalanon [[
                if(document.querySelector('.player-info__server-info-bar-multiplyer') != null)
                    document.querySelector('.player-info__server-info-bar-multiplyer').remove()
            ]]
        end
        if settings.deleteADDVIP[0] then
            evalanon([[
                document.querySelector('.player-info__vip').remove()
            ]])
        end
        evalanon(('window.executeEvent(\'event.arizonahud.serverInfo\',\'[%s]\');'):format(encodeJson(hud_data)))
    end
end

function setImages()
    local k = getweapon(getCurrentCharWeapon(PLAYER_PED))
    set(ffi.string(settings.weapons[allweapons[k].name]), false, '.player-info__gun-image', true)
    set(ffi.string(settings.flagurl), true, '.player-info__server-info-bar', false)
    set(ffi.string(settings.logourl), false, '.player-info__project-logo-image', true)
    set(ffi.string(settings.background), true, '.player-info__params', false)
    evalanon([[
        if(document.querySelector('.player-info__params') != null)
            document.querySelector('.player-info__params').style.backgroundSize = 'cover';
    ]])
    if settings.deleteBackground[0] then
        evalanon([[
            document.querySelector('.player-info__params').style.background = 'none'
        ]])
    end
    if settings.deleteWeaponBackground[0] then
        evalanon([[
            document.querySelector('.player-info__gun').style.background = 'none'
        ]])
    end
    local colorPlusMoney = string.format('rgba(%d, %d, %d, %f)', settings.colors.plusMoney[0]*255, settings.colors.plusMoney[1]*255, settings.colors.plusMoney[2]*255, settings.colors.plusMoney[3])
    local colorMinusMoney = string.format('rgba(%d, %d, %d, %f)', settings.colors.minusMoney[0]*255, settings.colors.minusMoney[1]*255, settings.colors.minusMoney[2]*255, settings.colors.minusMoney[3])
    evalanon(([[
        if(document.querySelector(`.player-info__cash--funding`) != null)
            document.querySelector(`.player-info__cash--funding`).style.color = '%s'
        if(document.querySelector(`.player-info__cash--withdrawl`) != null)
            document.querySelector(`.player-info__cash--withdrawl`).style.color = '%s'
    ]]):format(colorPlusMoney, colorMinusMoney))
    local colorGreenZoneText = string.format('rgba(%d, %d, %d, %f)', settings.colors.greenZoneText[0]*255, settings.colors.greenZoneText[1]*255, settings.colors.greenZoneText[2]*255, settings.colors.greenZoneText[3])
    local colorGreenZoneText2 = string.format('rgba(%d, %d, %d, %f)', settings.colors.greenZoneText2[0]*255, settings.colors.greenZoneText2[1]*255, settings.colors.greenZoneText2[2]*255, settings.colors.greenZoneText2[3])
    local colorGreenZoneBg = string.format('rgba(%d, %d, %d, %f)', settings.colors.greenZoneBg[0]*255, settings.colors.greenZoneBg[1]*255, settings.colors.greenZoneBg[2]*255, settings.colors.greenZoneBg[3])
    local colorGreenZoneGr1 = string.format('rgba(%d, %d, %d, %f)', settings.colors.greenZoneGr1[0]*255, settings.colors.greenZoneGr1[1]*255, settings.colors.greenZoneGr1[2]*255, settings.colors.greenZoneGr1[3])
    local colorGreenZoneGr2 = string.format('rgba(%d, %d, %d, %f)', settings.colors.greenZoneGr2[0]*255, settings.colors.greenZoneGr2[1]*255, settings.colors.greenZoneGr2[2]*255, settings.colors.greenZoneGr2[3])
    local colorGreenZoneIcon = string.format('rgba(%d, %d, %d, %f)', settings.colors.greenZoneIcon[0]*255, settings.colors.greenZoneIcon[1]*255, settings.colors.greenZoneIcon[2]*255, settings.colors.greenZoneIcon[3])
    evalanon(([[
            if (document.querySelector(`.player-info__green-zone-icon`) != null)
            {
                document.querySelector('.player-info__green-zone-text').style.color = "%s"
                document.querySelector('.player-info__green-zone-text-highlight').style.color = "%s"
                document.querySelector('.player-info__green-zone').style.background = "%s linear-gradient(90deg,%s,%s)"
                let greenzone = `<svg class="player-info__green-zone-icon" width="46" height="46" viewBox="0 0 46 46" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M41.4309 5.46996L23.3873 0.0568711C23.261 0.018957 23.1305 0 23 0V46H23.0001C23.1879 46 23.3756 45.9608 23.5508 45.8823C23.7401 45.7976 28.2411 43.7678 32.8121 40.4591C35.533 38.4894 37.7108 36.4722 39.2851 34.4633C41.3462 31.833 42.3912 29.1903 42.3912 26.6087V6.76074C42.3912 6.16562 42.0009 5.64093 41.4309 5.46996Z" fill="url(#paint0_linear_148_283)"/>
                <path d="M4.94096 5.46996L22.6205 0.0568711C22.7443 0.018957 22.8722 0 23 0V46H22.9999C22.8159 46 22.6319 45.9608 22.4603 45.8823C22.2748 45.7976 17.8646 43.7678 13.3859 40.4591C10.7198 38.4894 8.58598 36.4722 7.0435 34.4633C5.02398 31.833 4 29.1903 4 26.6087V6.76074C4 6.16562 4.38249 5.64093 4.94096 5.46996Z" fill="url(#paint1_linear_148_283)"/>
                <path d="M30.821 24.4903L29.8932 22.4068C29.7884 22.2405 29.6446 22.0994 29.4731 21.9943V25.6501C29.473 26.549 28.9034 27.3595 28.0296 27.7042C27.1558 28.0488 26.1497 27.8597 25.4797 27.225C25.4493 27.1962 25.4213 27.1673 25.3943 27.1385C25.0094 28.0437 24.0891 28.64 23.0616 28.65C22.64 28.6505 22.2261 28.5415 21.865 28.3347C21.865 28.3419 21.8637 28.3493 21.8623 28.3565C21.8611 28.3633 21.8598 28.37 21.8598 28.376C21.9151 29.2542 22.6575 29.9521 23.5823 29.9952C23.6466 29.9982 23.7109 29.9896 23.772 29.9703C24.3456 29.8112 24.9385 29.7242 25.5358 29.7116C26.4579 29.7547 27.3764 29.851 28.2861 30C29.0468 29.989 29.7118 29.5101 29.9288 28.8173L30.9516 25.4627C31.0447 25.145 31.0024 24.8051 30.8341 24.5166C30.8314 24.5118 30.8288 24.5069 30.8263 24.5019C30.8245 24.4981 30.8227 24.4943 30.821 24.4903Z" fill="black"/>
                <path fill-rule="evenodd" clip-rule="evenodd" d="M22.5 37C30.5081 37 37 30.5081 37 22.5C37 14.4919 30.5081 8 22.5 8C14.4919 8 8 14.4919 8 22.5C8 30.5081 14.4919 37 22.5 37ZM23.0616 28.0731C22.6283 28.0746 22.2113 27.9199 21.8944 27.6417L24.9287 23.4808V26.4231C24.8673 27.3575 24.0469 28.0826 23.0616 28.0731ZM16.5 33.4685L21.3327 26.7581L24.9287 21.765L25.5358 20.922L28.444 16.8838L31 13.3348C31.2656 13.5812 31.5205 13.8391 31.7639 14.1075C33.7748 16.326 35 19.2699 35 22.5C35 29.4036 29.4036 35 22.5 35C20.6598 35 18.9125 34.6024 17.3393 33.8883C17.0536 33.7586 16.7737 33.6186 16.5 33.4685ZM28.866 18.0814L25.5358 22.6482V25.6501C25.5358 26.2096 25.6509 26.5696 25.9088 26.8171C26.4052 27.2869 27.1503 27.4266 27.7973 27.1713C28.4442 26.9159 28.8659 26.3157 28.866 25.6501V18.0814ZM29.8932 22.4068L30.821 24.4903C30.8227 24.4943 30.8245 24.4981 30.8263 24.5019C30.8288 24.5069 30.8314 24.5118 30.8341 24.5166C31.0024 24.8051 31.0447 25.145 30.9516 25.4627L29.9288 28.8173C29.7118 29.5101 29.0468 29.989 28.2861 30C27.3764 29.851 26.4579 29.7547 25.5358 29.7116C24.9385 29.7242 24.3456 29.8112 23.772 29.9703C23.7109 29.9896 23.6466 29.9982 23.5823 29.9952C22.6575 29.9521 21.9151 29.2542 21.8598 28.376C21.8598 28.37 21.8611 28.3633 21.8623 28.3565C21.8637 28.3493 21.865 28.3419 21.865 28.3347C22.2261 28.5415 22.64 28.6505 23.0616 28.65C24.0891 28.64 25.0094 28.0437 25.3943 27.1385C25.4213 27.1673 25.4493 27.1962 25.4797 27.225C26.1497 27.8597 27.1558 28.0488 28.0296 27.7042C28.9034 27.3595 29.473 26.549 29.4731 25.6501V21.9943C29.6446 22.0994 29.7884 22.2405 29.8932 22.4068ZM13.5 31.1747C11.3328 28.9268 10 25.869 10 22.5C10 15.5964 15.5964 10 22.5 10C24.5411 10 26.468 10.4892 28.1697 11.3568C28.4529 11.5012 28.7298 11.656 29 11.8208L24.9287 17.5162L21.2858 22.6123L20.6787 23.4617L17.9137 27.3297L14.5 32.1051C14.1499 31.8132 13.8161 31.5026 13.5 31.1747Z" fill="white"/>
                <path d="M15.3266 26.9106C16.0914 26.8835 16.7045 26.3003 16.7321 25.5735V21.6346C16.7321 21.4753 16.8681 21.3461 17.0358 21.3461C17.2033 21.3461 17.3393 21.4753 17.3393 21.6346V25.9881L20.6787 21.4768V21.0577C20.6787 20.8984 20.8146 20.7692 20.9822 20.7692C21.0492 20.7692 21.1111 20.7898 21.1613 20.8247L24.8679 15.8173C24.6745 15.3048 24.0515 15 23.4684 15H22.6549C21.9225 15.0006 21.329 15.5646 21.3282 16.2605L21.2858 17.7164V17.8845C21.2858 18.0439 21.1499 18.1731 20.9822 18.1731C20.8146 18.1731 20.6787 18.0439 20.6787 17.8845V17.7087C20.6787 17.0827 20.3295 16.451 19.5494 16.451H18.584C18.2496 16.4363 17.9242 16.556 17.6874 16.7806C17.4507 17.0054 17.3242 17.3144 17.3393 17.6321V18.4615C17.3393 18.6208 17.2033 18.75 17.0358 18.75C16.8681 18.75 16.7321 18.6208 16.7321 18.4615V18.4081C16.7329 18.1761 16.6358 17.9536 16.4626 17.79C16.2503 17.6488 15.9917 17.5852 15.734 17.6106C14.7769 17.6112 14.001 18.3482 14 19.2576V25.6501C13.9992 25.9846 14.1387 26.3056 14.3877 26.5422C14.6366 26.7788 14.9746 26.9114 15.3266 26.9106Z" fill="white"/>
                <defs>
                <linearGradient id="paint0_linear_148_283" x1="32.6956" y1="0" x2="32.6956" y2="46" gradientUnits="userSpaceOnUse">
                <stop stop-color="%s"/>
                <stop offset="1" stop-color="%s" stop-opacity="0.15"/>
                </linearGradient>
                <linearGradient id="paint1_linear_148_283" x1="13.5" y1="0" x2="13.5" y2="46" gradientUnits="userSpaceOnUse">
                <stop stop-color="%s" stop-opacity="0.25"/>
                <stop offset="1" stop-color="%s"/>
                </linearGradient>
                </defs>
                </svg>
                `
                if (document.querySelector(`.player-info__green-zone-icon`)) document.querySelector(`.player-info__green-zone-icon`).remove()
                document.querySelector(`.player-info__green-zone`).insertAdjacentHTML("afterbegin", greenzone)
            }
        ]]):format(colorGreenZoneText, colorGreenZoneText2, colorGreenZoneBg, colorGreenZoneGr1, colorGreenZoneGr2, colorGreenZoneIcon, colorGreenZoneIcon, colorGreenZoneIcon, colorGreenZoneIcon))
end

function SoftBlueTheme()
    imgui.SwitchContext()
    local style = imgui.GetStyle()
  
    style.WindowPadding = imgui.ImVec2(15, 15)
    style.WindowRounding = 10.0
    style.ChildRounding = 6.0
    style.FramePadding = imgui.ImVec2(8, 7)
    style.FrameRounding = 8.0
    style.ItemSpacing = imgui.ImVec2(8, 8)
    style.ItemInnerSpacing = imgui.ImVec2(10, 6)
    style.IndentSpacing = 25.0
    style.ScrollbarSize = 13.0
    style.ScrollbarRounding = 12.0
    style.GrabMinSize = 10.0
    style.GrabRounding = 6.0
    style.PopupRounding = 8
    style.WindowTitleAlign = imgui.ImVec2(0.5, 0.5)
    style.ButtonTextAlign = imgui.ImVec2(0.5, 0.5)

    style.Colors[imgui.Col.Text]                   = imgui.ImVec4(0.90, 0.90, 0.93, 1.00)
    style.Colors[imgui.Col.TextDisabled]           = imgui.ImVec4(0.40, 0.40, 0.45, 1.00)
    style.Colors[imgui.Col.WindowBg]               = imgui.ImVec4(0.12, 0.12, 0.14, 1.00)
    style.Colors[imgui.Col.ChildBg]                = imgui.ImVec4(0.18, 0.20, 0.22, 0.30)
    style.Colors[imgui.Col.PopupBg]                = imgui.ImVec4(0.13, 0.13, 0.15, 1.00)
    style.Colors[imgui.Col.Border]                 = imgui.ImVec4(0.30, 0.30, 0.35, 1.00)
    style.Colors[imgui.Col.BorderShadow]           = imgui.ImVec4(0.00, 0.00, 0.00, 0.00)
    style.Colors[imgui.Col.FrameBg]                = imgui.ImVec4(0.18, 0.18, 0.20, 1.00)
    style.Colors[imgui.Col.FrameBgHovered]         = imgui.ImVec4(0.25, 0.25, 0.28, 1.00)
    style.Colors[imgui.Col.FrameBgActive]          = imgui.ImVec4(0.30, 0.30, 0.34, 1.00)
    style.Colors[imgui.Col.TitleBg]                = imgui.ImVec4(0.15, 0.15, 0.17, 1.00)
    style.Colors[imgui.Col.TitleBgCollapsed]       = imgui.ImVec4(0.10, 0.10, 0.12, 1.00)
    style.Colors[imgui.Col.TitleBgActive]          = imgui.ImVec4(0.15, 0.15, 0.17, 1.00)
    style.Colors[imgui.Col.MenuBarBg]              = imgui.ImVec4(0.12, 0.12, 0.14, 1.00)
    style.Colors[imgui.Col.ScrollbarBg]            = imgui.ImVec4(0.12, 0.12, 0.14, 1.00)
    style.Colors[imgui.Col.ScrollbarGrab]          = imgui.ImVec4(0.30, 0.30, 0.35, 1.00)
    style.Colors[imgui.Col.ScrollbarGrabHovered]   = imgui.ImVec4(0.40, 0.40, 0.45, 1.00)
    style.Colors[imgui.Col.ScrollbarGrabActive]    = imgui.ImVec4(0.50, 0.50, 0.55, 1.00)
    style.Colors[imgui.Col.CheckMark]              = imgui.ImVec4(0.70, 0.70, 0.90, 1.00)
    style.Colors[imgui.Col.SliderGrab]             = imgui.ImVec4(0.70, 0.70, 0.90, 1.00)
    style.Colors[imgui.Col.SliderGrabActive]       = imgui.ImVec4(0.80, 0.80, 0.90, 1.00)
    style.Colors[imgui.Col.Button]                 = imgui.ImVec4(0.18, 0.18, 0.20, 1.00)
    style.Colors[imgui.Col.ButtonHovered]          = imgui.ImVec4(0.60, 0.60, 0.90, 1.00)
    style.Colors[imgui.Col.ButtonActive]           = imgui.ImVec4(0.28, 0.56, 0.96, 1.00)
    style.Colors[imgui.Col.Header]                 = imgui.ImVec4(0.20, 0.20, 0.23, 1.00)
    style.Colors[imgui.Col.HeaderHovered]          = imgui.ImVec4(0.25, 0.25, 0.28, 1.00)
    style.Colors[imgui.Col.HeaderActive]           = imgui.ImVec4(0.30, 0.30, 0.34, 1.00)
    style.Colors[imgui.Col.Separator]              = imgui.ImVec4(0.40, 0.40, 0.45, 1.00)
    style.Colors[imgui.Col.SeparatorHovered]       = imgui.ImVec4(0.50, 0.50, 0.55, 1.00)
    style.Colors[imgui.Col.SeparatorActive]        = imgui.ImVec4(0.60, 0.60, 0.65, 1.00)
    style.Colors[imgui.Col.ResizeGrip]             = imgui.ImVec4(0.20, 0.20, 0.23, 1.00)
    style.Colors[imgui.Col.ResizeGripHovered]      = imgui.ImVec4(0.25, 0.25, 0.28, 1.00)
    style.Colors[imgui.Col.ResizeGripActive]       = imgui.ImVec4(0.30, 0.30, 0.34, 1.00)
    style.Colors[imgui.Col.PlotLines]              = imgui.ImVec4(0.61, 0.61, 0.64, 1.00)
    style.Colors[imgui.Col.PlotLinesHovered]       = imgui.ImVec4(0.70, 0.70, 0.75, 1.00)
    style.Colors[imgui.Col.PlotHistogram]          = imgui.ImVec4(0.61, 0.61, 0.64, 1.00)
    style.Colors[imgui.Col.PlotHistogramHovered]   = imgui.ImVec4(0.70, 0.70, 0.75, 1.00)
    style.Colors[imgui.Col.TextSelectedBg]         = imgui.ImVec4(0.30, 0.30, 0.34, 1.00)
    style.Colors[imgui.Col.ModalWindowDimBg]       = imgui.ImVec4(0.10, 0.10, 0.12, 0.80)
    style.Colors[imgui.Col.Tab]                    = imgui.ImVec4(0.18, 0.20, 0.22, 1.00)
    style.Colors[imgui.Col.TabHovered]             = imgui.ImVec4(0.60, 0.60, 0.90, 1.00)
    style.Colors[imgui.Col.TabActive]              = imgui.ImVec4(0.28, 0.56, 0.96, 1.00)
end

function getweapon(wid)
    for k, v in pairs(allweapons) do
        if v.id == wid then
            return k
        end
    end
    return 0
end

function getAllConfigs()
    local files = {}
    for k, file in pairs(get_all_files_info_async(script_path)) do
        if file:find('.+%.json') and file ~= "default.json" then
            table.insert(files, file)
        end
    end
    return files
end

function get_all_files_info_async(path)
    local entries = {}
    for entry in lfs.dir(path) do
        table.insert(entries, entry)
    end
    return entries
end

function check(var) local result if var == nil then return false end if type(var) == "string" then if #var > 0 then result = var end elseif type(var) == "number" then if var > 0 then result = var end end return result end --говно