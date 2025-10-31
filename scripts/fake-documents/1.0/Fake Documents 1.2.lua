script_name("Fake Documents")
script_author("DepsCian")
script_version("1.0")

local imgui = require("mimgui")
local encoding = require("encoding")
local ffi = require("ffi")
local acef =require("arizona-events")

encoding.default = "CP1251"
local u8 = encoding.UTF8

local json_config_dir = getWorkingDirectory() .. "\\\\config"
local json_config_path = json_config_dir .. "\\\\fake_doc.json"

local notlocal_state = false

local default_config = {
    pass = {
        enabled         = false,
        onlyOwn         = true,
        name            = { value = u8"Олег Тиньков", enabled = true },
        sex             = { value = u8"Мужской", enabled = true },
        birthday        = { value = u8"25.12.1967", enabled = true },
        citizen         = { value = u8"Российская Федерация", enabled = true },
        married         = { value = u8"Женат", enabled = true },
        level           = { value = u8"99", enabled = true },
        zakono          = { value = u8"777", enabled = true },
        job             = { value = u8"Предприниматель", enabled = true },
        povestka        = { value = u8"Нет", enabled = true },
        rank            = { value = u8"CEO", enabled = true },
        organization    = { value = u8"Tinkoff Bank", enabled = true },
        seria           = { value = u8"0011", enabled = true },
        number          = { value = u8"777777", enabled = true },
        signature       = { value = u8"O.Tinkoff", enabled = true },
        avatarUrl       = { value = u8"https://storage.googleapis.com/activatica/uploads/f39fc1d4-16df-4d3e-97de-5eaa235edd3c", enabled = true }
    },
    medical = {
        enabled         = false,
        onlyOwn         = true,
        name            = { value = u8"Александр Невский", enabled = true },
        ukrop           = { value = u8"0.0", enabled = true },
        health          = { value = u8"Полностью здоров", enabled = true },
        insurance       = { value = u8"Platinum", enabled = true },
        validity        = { value = u8"999 дня(ей)", enabled = true },
        psychiatric     = { value = u8"0 раз(а)", enabled = true },
        examProgress    = false,
        avatarUrl       = { value = u8"", enabled = false }
    },
    licenses = {
        enabled     = false,
        onlyOwn     = true,
        car         = { enabled = true, date = u8"12:12 11.05.2029" },
        bike        = { enabled = true, date = u8"12:12 11.05.2029" },
        airship     = { enabled = true, date = u8"12:12 11.05.2029" },
        boat        = { enabled = true, date = u8"12:12 11.05.2029" },
        fishing     = { enabled = true, date = u8"12:12 11.05.2029" },
        glock       = { enabled = true, date = u8"12:12 11.05.2029" },
        deer        = { enabled = true, date = u8"12:12 11.05.2029" },
        shovel      = { enabled = true, date = u8"12:12 11.05.2029" },
        taxi        = { enabled = true, date = u8"12:12 11.05.2029" },
        repair      = { enabled = true, date = u8"12:12 11.05.2029" },
        court       = { enabled = true, date = u8"12:12 11.05.2029" },
        tax         = { enabled = true, date = u8"12:12 11.05.2029" },
        diplomacy   = { enabled = true, date = u8"12:12 11.05.2029" },
        pickaxe     = { enabled = true, date = u8"12:12 11.05.2029" }
    }
}

local function createConfigDirectoryIfNotExists()
    if not doesDirectoryExist(json_config_dir) then createDirectory(json_config_dir) end
end

local function jsonSave(data, path)
    local file = io.open(path, "w")
    if file then
        file:write(encodeJson(data))
        file:flush()
        file:close()
        return true
    end
    return false
end

local function jsonRead(path)
    local file = io.open(path, "r")
    if file then
        local content = file:read("*a")
        file:close()
        local success, result = pcall(decodeJson, content)
        if success then return result end
    end
    return nil
end

createConfigDirectoryIfNotExists()
local cfg = jsonRead(json_config_path)
if not cfg then
    cfg = default_config
    jsonSave(cfg, json_config_path)
end

local window = imgui.new.bool(false)
local tab = imgui.new.int(0)

local passData = {
    enabled = imgui.new.bool(cfg.pass.enabled),
    onlyOwn = imgui.new.bool(cfg.pass.onlyOwn ~= false),
    name = {value = imgui.new.char[256](cfg.pass.name.value), enabled = imgui.new.bool(cfg.pass.name.enabled)},
    sex = {value = imgui.new.char[64](cfg.pass.sex.value), enabled = imgui.new.bool(cfg.pass.sex.enabled)},
    birthday = {value = imgui.new.char[64](cfg.pass.birthday.value), enabled = imgui.new.bool(cfg.pass.birthday.enabled)},
    citizen = {value = imgui.new.char[128](cfg.pass.citizen.value), enabled = imgui.new.bool(cfg.pass.citizen.enabled)},
    married = {value = imgui.new.char[64](cfg.pass.married.value), enabled = imgui.new.bool(cfg.pass.married.enabled)},
    level = {value = imgui.new.char[64](cfg.pass.level.value), enabled = imgui.new.bool(cfg.pass.level.enabled)},
    zakono = {value = imgui.new.char[64](cfg.pass.zakono.value), enabled = imgui.new.bool(cfg.pass.zakono.enabled)},
    job = {value = imgui.new.char[128](cfg.pass.job.value), enabled = imgui.new.bool(cfg.pass.job.enabled)},
    povestka = {value = imgui.new.char[128](cfg.pass.povestka.value), enabled = imgui.new.bool(cfg.pass.povestka.enabled)},
    rank = {value = imgui.new.char[128](cfg.pass.rank.value), enabled = imgui.new.bool(cfg.pass.rank.enabled)},
    organization = {value = imgui.new.char[128](cfg.pass.organization.value), enabled = imgui.new.bool(cfg.pass.organization.enabled)},
    seria = {value = imgui.new.char[32](cfg.pass.seria.value), enabled = imgui.new.bool(cfg.pass.seria.enabled)},
    number = {value = imgui.new.char[32](cfg.pass.number.value), enabled = imgui.new.bool(cfg.pass.number.enabled)},
    signature = {value = imgui.new.char[128](cfg.pass.signature.value), enabled = imgui.new.bool(cfg.pass.signature.enabled)},
    avatarUrl = {value = imgui.new.char[512](cfg.pass.avatarUrl.value), enabled = imgui.new.bool(cfg.pass.avatarUrl.enabled)}
}

local licenseData = {
    enabled = imgui.new.bool(cfg.licenses.enabled),
    onlyOwn = imgui.new.bool(cfg.licenses.onlyOwn ~= false),
    {name = u8"Лицензия на авто", id = "car", icon = "icon-car", enabled = imgui.new.bool(cfg.licenses.car.enabled), date = imgui.new.char[64](cfg.licenses.car.date)},
    {name = u8"Лицензия на мото", id = "bike", icon = "icon-bike", enabled = imgui.new.bool(cfg.licenses.bike.enabled), date = imgui.new.char[64](cfg.licenses.bike.date)},
    {name = u8"Лицензия на полеты", id = "airship", icon = "icon-airship", enabled = imgui.new.bool(cfg.licenses.airship.enabled), date = imgui.new.char[64](cfg.licenses.airship.date)},
    {name = u8"Лицензия на плавание", id = "boat", icon = "icon-boat", enabled = imgui.new.bool(cfg.licenses.boat.enabled), date = imgui.new.char[64](cfg.licenses.boat.date)},
    {name = u8"Лицензия на ловлю рыбы", id = "fishing", icon = "icon-fishing", enabled = imgui.new.bool(cfg.licenses.fishing.enabled), date = imgui.new.char[64](cfg.licenses.fishing.date)},
    {name = u8"Лицензия на оружие", id = "glock", icon = "icon-glock", enabled = imgui.new.bool(cfg.licenses.glock.enabled), date = imgui.new.char[64](cfg.licenses.glock.date)},
    {name = u8"Лицензия на охоту", id = "deer", icon = "icon-deer", enabled = imgui.new.bool(cfg.licenses.deer.enabled), date = imgui.new.char[64](cfg.licenses.deer.date)},
    {name = u8"Лицензия на раскопки", id = "shovel", icon = "icon-shovel", enabled = imgui.new.bool(cfg.licenses.shovel.enabled), date = imgui.new.char[64](cfg.licenses.shovel.date)},
    {name = u8"Лицензия таксиста", id = "taxi", icon = "icon-taxi", enabled = imgui.new.bool(cfg.licenses.taxi.enabled), date = imgui.new.char[64](cfg.licenses.taxi.date)},
    {name = u8"Лицензия механика", id = "repair", icon = "icon-repair", enabled = imgui.new.bool(cfg.licenses.repair.enabled), date = imgui.new.char[64](cfg.licenses.repair.date)},
    {name = u8"Лицензия адвоката", id = "court", icon = "icon-court", enabled = imgui.new.bool(cfg.licenses.court.enabled), date = imgui.new.char[64](cfg.licenses.court.date)},
    {name = u8"Лицензия налоговика", id = "tax", icon = "icon-tax", enabled = imgui.new.bool(cfg.licenses.tax.enabled), date = imgui.new.char[64](cfg.licenses.tax.date)},
    {name = u8"Лицензия дипломата", id = "diplomacy", icon = "icon-diplomacy", enabled = imgui.new.bool(cfg.licenses.diplomacy.enabled), date = imgui.new.char[64](cfg.licenses.diplomacy.date)},
    {name = u8"Разрешение на добычу ресурсов", id = "pickaxe", icon = "icon-pickaxe", enabled = imgui.new.bool(cfg.licenses.pickaxe.enabled), date = imgui.new.char[64](cfg.licenses.pickaxe.date)}
}

local medicalData = {
    enabled = imgui.new.bool(cfg.medical.enabled),
    onlyOwn = imgui.new.bool(cfg.medical.onlyOwn ~= false),
    name = {value = imgui.new.char[256](cfg.medical.name.value), enabled = imgui.new.bool(cfg.medical.name.enabled)},
    ukrop = {value = imgui.new.char[32](cfg.medical.ukrop.value), enabled = imgui.new.bool(cfg.medical.ukrop.enabled)},
    health = {value = imgui.new.char[128](cfg.medical.health.value), enabled = imgui.new.bool(cfg.medical.health.enabled)},
    insurance = {value = imgui.new.char[128](cfg.medical.insurance.value), enabled = imgui.new.bool(cfg.medical.insurance.enabled)},
    validity = {value = imgui.new.char[128](cfg.medical.validity.value), enabled = imgui.new.bool(cfg.medical.validity.enabled)},
    psychiatric = {value = imgui.new.char[64](cfg.medical.psychiatric.value), enabled = imgui.new.bool(cfg.medical.psychiatric.enabled)},
    examProgress = imgui.new.bool(cfg.medical.examProgress),
    avatarUrl = {value = imgui.new.char[512](cfg.medical.avatarUrl.value), enabled = imgui.new.bool(cfg.medical.avatarUrl.enabled)}
}

function saveConfig()
    cfg.pass.enabled = passData.enabled[0]
    cfg.pass.onlyOwn = passData.onlyOwn[0]
    cfg.pass.name.value = ffi.string(passData.name.value)
    cfg.pass.name.enabled = passData.name.enabled[0]
    cfg.pass.sex.value = ffi.string(passData.sex.value)
    cfg.pass.sex.enabled = passData.sex.enabled[0]
    cfg.pass.birthday.value = ffi.string(passData.birthday.value)
    cfg.pass.birthday.enabled = passData.birthday.enabled[0]
    cfg.pass.citizen.value = ffi.string(passData.citizen.value)
    cfg.pass.citizen.enabled = passData.citizen.enabled[0]
    cfg.pass.married.value = ffi.string(passData.married.value)
    cfg.pass.married.enabled = passData.married.enabled[0]
    cfg.pass.level.value = ffi.string(passData.level.value)
    cfg.pass.level.enabled = passData.level.enabled[0]
    cfg.pass.zakono.value = ffi.string(passData.zakono.value)
    cfg.pass.zakono.enabled = passData.zakono.enabled[0]
    cfg.pass.job.value = ffi.string(passData.job.value)
    cfg.pass.job.enabled = passData.job.enabled[0]
    cfg.pass.povestka.value = ffi.string(passData.povestka.value)
    cfg.pass.povestka.enabled = passData.povestka.enabled[0]
    cfg.pass.rank.value = ffi.string(passData.rank.value)
    cfg.pass.rank.enabled = passData.rank.enabled[0]
    cfg.pass.organization.value = ffi.string(passData.organization.value)
    cfg.pass.organization.enabled = passData.organization.enabled[0]
    cfg.pass.seria.value = ffi.string(passData.seria.value)
    cfg.pass.seria.enabled = passData.seria.enabled[0]
    cfg.pass.number.value = ffi.string(passData.number.value)
    cfg.pass.number.enabled = passData.number.enabled[0]
    cfg.pass.signature.value = ffi.string(passData.signature.value)
    cfg.pass.signature.enabled = passData.signature.enabled[0]
    cfg.pass.avatarUrl.value = ffi.string(passData.avatarUrl.value)
    cfg.pass.avatarUrl.enabled = passData.avatarUrl.enabled[0]
    
    cfg.medical.enabled = medicalData.enabled[0]
    cfg.medical.onlyOwn = medicalData.onlyOwn[0]
    cfg.medical.name.value = ffi.string(medicalData.name.value)
    cfg.medical.name.enabled = medicalData.name.enabled[0]
    cfg.medical.ukrop.value = ffi.string(medicalData.ukrop.value)
    cfg.medical.ukrop.enabled = medicalData.ukrop.enabled[0]
    cfg.medical.health.value = ffi.string(medicalData.health.value)
    cfg.medical.health.enabled = medicalData.health.enabled[0]
    cfg.medical.insurance.value = ffi.string(medicalData.insurance.value)
    cfg.medical.insurance.enabled = medicalData.insurance.enabled[0]
    cfg.medical.validity.value = ffi.string(medicalData.validity.value)
    cfg.medical.validity.enabled = medicalData.validity.enabled[0]
    cfg.medical.psychiatric.value = ffi.string(medicalData.psychiatric.value)
    cfg.medical.psychiatric.enabled = medicalData.psychiatric.enabled[0]
    cfg.medical.examProgress = medicalData.examProgress[0]
    cfg.medical.avatarUrl.value = ffi.string(medicalData.avatarUrl.value)
    cfg.medical.avatarUrl.enabled = medicalData.avatarUrl.enabled[0]
    
    cfg.licenses.enabled = licenseData.enabled[0]
    cfg.licenses.onlyOwn = licenseData.onlyOwn[0]
    for i, license in ipairs(licenseData) do
        if i > 0 and license.id then
            cfg.licenses[license.id].enabled = license.enabled[0]
            cfg.licenses[license.id].date = ffi.string(license.date)
        end
    end
    
    jsonSave(cfg, json_config_path)
end

function main()
	if not isSampLoaded() or not isSampfuncsLoaded() then return end
	while not isSampAvailable() do wait(1000) end

    sampRegisterChatCommand("fdoc", function() window[0] = not window[0] end)

    while true do wait(0) end
end

imgui.OnFrame(
    function() return window[0] end,
    function()
        local width, height = 600, 400
        local screenWidth, screenHeight = getScreenResolution()
        imgui.SetNextWindowPos(imgui.ImVec2(screenWidth / 2, screenHeight / 2), imgui.Cond.FirstUseEver, imgui.ImVec2(0.5, 0.5))
        imgui.SetNextWindowSize(imgui.ImVec2(width, height), imgui.Cond.FirstUseEver)
        
        if imgui.Begin(u8"Fake Documents by DepsCian", window) then
            imgui.BeginTabBar("TabBar")
            
            if imgui.BeginTabItem(u8"Паспорт") then
                tab[0] = 0
                passTab()
                imgui.EndTabItem()
            end
            
            if imgui.BeginTabItem(u8"Лицензии") then
                tab[0] = 1
                licensesTab()
                imgui.EndTabItem()
            end
            
            if imgui.BeginTabItem(u8"Мед. карта") then
                tab[0] = 2
                medicalTab()
                imgui.EndTabItem()
            end
            
            imgui.EndTabBar()
            
            if imgui.Button(u8"Сохранить настройки") then
                saveConfig()
                sampAddChatMessage("{00AAFF}[Fake Documents]: {FFFFFF}Настройки успешно сохранены", -1)
            end
        end
        imgui.End()
    end
)

function renderRow(label, id, setting)
    imgui.Text(label)
    imgui.NextColumn()
    imgui.PushItemWidth(imgui.GetColumnWidth() - 10)
    imgui.InputText("##" .. id, setting.value, ffi.sizeof(setting.value))
    imgui.PopItemWidth()
    imgui.NextColumn()
    imgui.Checkbox("##" .. id .. "_enabled", setting.enabled)
    imgui.NextColumn()
end

function passTab()
    if imgui.Checkbox(u8"Включено", passData.enabled) then
        saveConfig()
    end
    imgui.SameLine()
    if imgui.Checkbox(u8"Не менять чужие документы", passData.onlyOwn) then
        saveConfig()
    end
    imgui.Separator()
    
    imgui.Columns(3)
    imgui.Text(u8"Параметр")
    imgui.NextColumn()
    imgui.Text(u8"Значение")
    imgui.NextColumn()
    imgui.Text(u8"Включено")
    imgui.NextColumn()
    imgui.Separator()
    
    renderRow(u8"Имя и фамилия", "name", passData.name)
    renderRow(u8"Пол", "sex", passData.sex)
    renderRow(u8"Дата рождения", "birthday", passData.birthday)
    renderRow(u8"Гражданство", "citizen", passData.citizen)
    renderRow(u8"Семейное положение", "married", passData.married)
    renderRow(u8"Уровень", "level", passData.level)
    renderRow(u8"Законопослушность", "zakono", passData.zakono)
    renderRow(u8"Работа", "job", passData.job)
    renderRow(u8"Повестка", "povestka", passData.povestka)
    renderRow(u8"Ранг", "rank", passData.rank)
    renderRow(u8"Организация", "organization", passData.organization)
    renderRow(u8"Серия", "seria", passData.seria)
    renderRow(u8"Номер", "number", passData.number)
    renderRow(u8"Подпись", "signature", passData.signature)
    renderRow(u8"URL фото", "avatarUrl", passData.avatarUrl)
    
    imgui.Columns(1)
end

function licensesTab()
    if imgui.Checkbox(u8"Включено", licenseData.enabled) then
        saveConfig()
    end
    imgui.SameLine()
    if imgui.Checkbox(u8"Не менять чужие документы", licenseData.onlyOwn) then
        saveConfig()
    end
    imgui.Separator()
    
    imgui.Columns(3)
    imgui.Text(u8"Название")
    imgui.NextColumn()
    imgui.Text(u8"Активна")
    imgui.NextColumn()
    imgui.Text(u8"Срок действия")
    imgui.NextColumn()
    imgui.Separator()
    
    for i, license in ipairs(licenseData) do
        imgui.Text(license.name)
        imgui.NextColumn()
        if imgui.Checkbox("##" .. i, license.enabled) then
            saveConfig()
        end
        imgui.NextColumn()
        imgui.PushItemWidth(imgui.GetColumnWidth() - 10)
        imgui.InputText("##date" .. i, license.date, ffi.sizeof(license.date))
        imgui.PopItemWidth()
        imgui.NextColumn()
    end
    
    imgui.Columns(1)
end

function medicalTab()
    if imgui.Checkbox(u8"Включено", medicalData.enabled) then
        saveConfig()
    end
    imgui.SameLine()
    if imgui.Checkbox(u8"Не менять чужие документы", medicalData.onlyOwn) then
        saveConfig()
    end
    imgui.Separator()
    
    imgui.Columns(3)
    imgui.Text(u8"Параметр")
    imgui.NextColumn()
    imgui.Text(u8"Значение")
    imgui.NextColumn()
    imgui.Text(u8"Включено")
    imgui.NextColumn()
    imgui.Separator()
    
    renderRow(u8"Имя и фамилия", "med_name", medicalData.name)
    renderRow(u8"Зависимость от укропа", "med_ukrop", medicalData.ukrop)
    renderRow(u8"Состояние здоровья", "med_health", medicalData.health)
    renderRow(u8"Медицинская страховка", "med_insurance", medicalData.insurance)
    renderRow(u8"Срок действия", "med_validity", medicalData.validity)
    renderRow(u8"Лечился в псих.больнице", "med_psychiatric", medicalData.psychiatric)
    renderRow(u8"URL фото", "med_avatarUrl", medicalData.avatarUrl)
    
    imgui.Columns(1)
    
    imgui.Separator()
    imgui.Checkbox(u8"Медосмотр пройден (10/10)", medicalData.examProgress)
end

-- window.executeEvent('event.documents.inititalizeData', `[{"type":1,"name":"Zerno_BamBas","sex":"Мужской","birthday":"06.11.2006","citizen":"Имеется (с рождения)","married":"Не женат(а)","level":"36 лет","zakono":"100/100","job":"Дальнобойщик","agenda":"Нет","charity":"TV студия SF","rank":"Журналист","seria":"8275","number":"588015","signature":"ZBamBas","skin_image_url":"https://cdn.azresources.cloud/projects/arizona-rp/assets/images/inventory/skins/512/947.webp"}]`); | 220, 17, 0, 0, 0, 0, 217, 1, 1, 119, 105, 110, 100, 111, 119, 46, 101, 120, 101, 99, 117, 116, 101, 69, 118, 101, 110, 116, 40, 39, 101, 118, 101, 110, 116, 46, 100, 111, 99, 117, 109, 101, 110, 116, 115, 46, 105, 110, 105, 116, 105, 116, 97, 108, 105, 122, 101, 68, 97, 116, 97, 39, 44, 32, 96, 91, 123, 34, 116, 121, 112, 101, 34, 58, 49, 44, 34, 110, 97, 109, 101, 34, 58, 34, 90, 101, 114, 110, 111, 95, 66, 97, 109, 66, 97, 115, 34, 44, 34, 115, 101, 120, 34, 58, 34, 204, 243, 230, 241, 234, 238, 233, 34, 44, 34, 98, 105, 114, 116, 104, 100, 97, 121, 34, 58, 34, 48, 54, 46, 49, 49, 46, 50, 48, 48, 54, 34, 44, 34, 99, 105, 116, 105, 122, 101, 110, 34, 58, 34, 200, 236, 229, 229, 242, 241, 255, 32, 40, 241, 32, 240, 238, 230, 228, 229, 237, 232, 255, 41, 34, 44, 34, 109, 97, 114, 114, 105, 101, 100, 34, 58, 34, 205, 229, 32, 230, 229, 237, 224, 242, 40, 224, 41, 34, 44, 34, 108, 101, 118, 101, 108, 34, 58, 34, 51, 54, 32, 235, 229, 242, 34, 44, 34, 122, 97, 107, 111, 110, 111, 34, 58, 34, 49, 48, 48, 47, 49, 48, 48, 34, 44, 34, 106, 111, 98, 34, 58, 34, 196, 224, 235, 252, 237, 238, 225, 238, 233, 249, 232, 234, 34, 44, 34, 97, 103, 101, 110, 100, 97, 34, 58, 34, 205, 229, 242, 34, 44, 34, 99, 104, 97, 114, 105, 116, 121, 34, 58, 34, 84, 86, 32, 241, 242, 243, 228, 232, 255, 32, 83, 70, 34, 44, 34, 114, 97, 110, 107, 34, 58, 34, 198, 243, 240, 237, 224, 235, 232, 241, 242, 34, 44, 34, 115, 101, 114, 105, 97, 34, 58, 34, 56, 50, 55, 53, 34, 44, 34, 110, 117, 109, 98, 101, 114, 34, 58, 34, 53, 56, 56, 48, 49, 53, 34, 44, 34, 115, 105, 103, 110, 97, 116, 117, 114, 101, 34, 58, 34, 90, 66, 97, 109, 66, 97, 115, 34, 44, 34, 115, 107, 105, 110, 95, 105, 109, 97, 103, 101, 95, 117, 114, 108, 34, 58, 34, 104, 116, 116, 112, 115, 58, 47, 47, 99, 100, 110, 46, 97, 122, 114, 101, 115, 111, 117, 114, 99, 101, 115, 46, 99, 108, 111, 117, 100, 47, 112, 114, 111, 106, 101, 99, 116, 115, 47, 97, 114, 105, 122, 111, 110, 97, 45, 114, 112, 47, 97, 115, 115, 101, 116, 115, 47, 105, 109, 97, 103, 101, 115, 47, 105, 110, 118, 101, 110, 116, 111, 114, 121, 47, 115, 107, 105, 110, 115, 47, 53, 49, 50, 47, 57, 52, 55, 46, 119, 101, 98, 112, 34, 125, 93, 96, 41, 59
-- razrabi arizoni, nahui vi mne golovu lomaete? bezgramotniye pizdec, NE inititalizeData, A initializeData!!!!!

function acef.onArizonaDisplay(packet)
    if not packet.text:find(sampGetPlayerNickname(select(2, sampGetPlayerIdByCharHandle(PLAYER_PED)))) and packet.text:find("event.documents.inititalizeData") then
        notlocal_state = true
    end
    -- Passport
    if packet.text:find('"type":1') or packet.text == "window.executeEvent('event.documents.updatePage', `[1]`);" then
        if passData.enabled[0] and (not passData.onlyOwn[0] or not notlocal_state) then evalanon(buildPassport()) end
    end
    -- Licenses
    if packet.text:find('"type":2') or packet.text == "window.executeEvent('event.documents.updatePage', `[2]`);" then
        if licenseData.enabled[0] and (not licenseData.onlyOwn[0] or not notlocal_state) then evalanon(buildLicenses()) end
    end
    -- Medical Card
    if packet.text:find('"type":4') or packet.text == "window.executeEvent('event.documents.updatePage', `[4]`);" then
        if medicalData.enabled[0] and (not medicalData.onlyOwn[0] or not notlocal_state) then evalanon(buildMedical()) end
    end
end

function acef.onArizonaSend(packet)
    if packet.text:find("onActiveViewChanged|null") then
        notlocal_state = false
    end
end

function buildPassport()
    local code = ""
    local selectors = {
        {name = "rank", sel = "body > div.documents > div.documents__content.documents__content--pasport > div > div > div.documents-pasport__main-info.svelte-11hrewx > div:nth-child(10) > div.documents-pasport__rank.svelte-11hrewx"},
        {name = "avatarUrl", sel = "body > div.documents > div.documents__content.documents__content--pasport > div > div > div.documents-pasport__photo-wrapper.svelte-11hrewx > img.documents-pasport__photo.svelte-11hrewx", attr = "src", style = "height:100%!important"},
        {name = "name", sel = "body > div.documents > div.documents__content.documents__content--pasport > div > div > div.documents-pasport__main-info.svelte-11hrewx > div:nth-child(1) > div.documents-pasport__value.svelte-11hrewx"},
        {name = "level", sel = "body > div.documents > div.documents__content.documents__content--pasport > div > div > div.documents-pasport__main-info.svelte-11hrewx > div:nth-child(2) > div.documents-pasport__value.svelte-11hrewx"},
        {name = "sex", sel = "body > div.documents > div.documents__content.documents__content--pasport > div > div > div.documents-pasport__main-info.svelte-11hrewx > div:nth-child(3) > div.documents-pasport__value.svelte-11hrewx"},
        {name = "zakono", sel = "body > div.documents > div.documents__content.documents__content--pasport > div > div > div.documents-pasport__main-info.svelte-11hrewx > div:nth-child(4) > div.documents-pasport__value.svelte-11hrewx"},
        {name = "birthday", sel = "body > div.documents > div.documents__content.documents__content--pasport > div > div > div.documents-pasport__main-info.svelte-11hrewx > div:nth-child(5) > div.documents-pasport__value.svelte-11hrewx"},
        {name = "job", sel = "body > div.documents > div.documents__content.documents__content--pasport > div > div > div.documents-pasport__main-info.svelte-11hrewx > div:nth-child(6) > div.documents-pasport__value.svelte-11hrewx"},
        {name = "citizen", sel = "body > div.documents > div.documents__content.documents__content--pasport > div > div > div.documents-pasport__main-info.svelte-11hrewx > div:nth-child(7) > div.documents-pasport__value.svelte-11hrewx"},
        {name = "povestka", sel = "body > div.documents > div.documents__content.documents__content--pasport > div > div > div.documents-pasport__main-info.svelte-11hrewx > div:nth-child(8) > div.documents-pasport__value.svelte-11hrewx"},
        {name = "married", sel = "body > div.documents > div.documents__content.documents__content--pasport > div > div > div.documents-pasport__main-info.svelte-11hrewx > div:nth-child(9) > div.documents-pasport__value.svelte-11hrewx"},
        {name = "organization", sel = "body > div.documents > div.documents__content.documents__content--pasport > div > div > div.documents-pasport__main-info.svelte-11hrewx > div:nth-child(10) > div.documents-pasport__value.svelte-11hrewx"},
        {name = "seria", sel = "body > div.documents > div.documents__content.documents__content--pasport > div > div > div.documents-pasport__serial.svelte-11hrewx > div.documents-pasport__serial-value.svelte-11hrewx"},
        {name = "number", sel = "body > div.documents > div.documents__content.documents__content--pasport > div > div > div.documents-pasport__number.svelte-11hrewx > div.documents-pasport__number-value.svelte-11hrewx"},
        {name = "signature", sel = "body > div.documents > div.documents__content.documents__content--pasport > div > div > div.documents-pasport__signature.svelte-11hrewx"}
    }
    
    for _, item in ipairs(selectors) do
        local setting = passData[item.name]
        if setting and setting.enabled[0] then
            if item.attr then
                code = code .. "document.querySelector(\"" .. item.sel .. "\")." .. item.attr .. "=\"" .. u8:decode(ffi.string(setting.value)) .. "\";"
                if item.style then
                    code = code .. "document.querySelector(\"" .. item.sel .. "\").style=\"" .. item.style .. "\";"
                end
            else
                code = code .. "document.querySelector(\"" .. item.sel .. "\").textContent=\"" .. u8:decode(ffi.string(setting.value)) .. "\";"
            end
        end
    end
    return code
end

function buildMedical()
    local code = ""
    if medicalData.avatarUrl.enabled[0] then
        code = code .. "document.querySelector(\".documents-medical__photo\").src=\"" .. u8:decode(ffi.string(medicalData.avatarUrl.value)) .. "\";"
        code = code .. "document.querySelector(\".documents-medical__photo\").style=\"height:100%!important\";"
    end
    local medSelectors = {
        {name = "name", sel = ".documents-medical__main-info > div:nth-child(1) > div.documents-medical__value"},
        {name = "ukrop", sel = ".documents-medical__main-info > div:nth-child(2) > div.documents-medical__value"},
        {name = "health", sel = ".documents-medical__main-info > div:nth-child(3) > div.documents-medical__value"},
        {name = "insurance", sel = ".documents-medical__main-info > div:nth-child(4) > div.documents-medical__value"},
        {name = "validity", sel = ".documents-medical__main-info > div:nth-child(5) > div.documents-medical__value"},
        {name = "psychiatric", sel = ".documents-medical__main-info > div:nth-child(6) > div.documents-medical__value-wrapper > div.documents-medical__value"},
    }
    for _, item in ipairs(medSelectors) do
        if medicalData[item.name].enabled[0] then
            code = code .. "document.querySelector(\"" .. item.sel .. "\").textContent=\"" .. u8:decode(ffi.string(medicalData[item.name].value)) .. "\";"
        end
    end
    if medicalData.psychiatric.enabled[0] then
        code = code .. "const r=document.querySelector(\".documents-medical__refresh-card\");if(r)r.style.display=\"none\";"
    end
    if medicalData.examProgress[0] then
        code = code .. "document.querySelector(\".documents-medical__inspection-progress\").textContent=\"10 / 10\";const p=document.querySelectorAll(\".documents-medical__progress-counter-item\");p.forEach(i=>{i.style.backgroundColor=\"#4CAF50\"});document.querySelector(\".documents-medical__tip-text\").textContent=\"Медицинский осмотр полностью пройден\";"
    end
    return code
end

function buildLicenses()
    local json = "{"
    for i, license in ipairs(licenseData) do
        local enabled = license.enabled[0] and "true" or "false"
        json = json .. '"' .. license.id .. '": {"enabled": ' .. enabled .. ', "date": "' .. u8:decode(ffi.string(license.date)) .. '"}'
        if i < #licenseData then
            json = json .. ", "
        end
    end
    json = json .. "}"
    local code = [[
        const l=document.querySelectorAll("body > div.documents > div.documents__content.documents__content--license > div > div > div.documents-license__license-grid.svelte-e82vsm > div.documents-license__license.svelte-e82vsm");
        const s=]] .. json .. [[;l.forEach(c=>{const i=c.querySelector("i.documents-license__license-icon");if(!i)return;let t=null;for(const c of i.classList){if(c.startsWith("icon-")){t=c.replace("icon-","");break}}if(!t||!s[t])return;const n=c.querySelector(".documents-license__license-icon"),a=c.querySelector(".documents-license__license-name"),e=c.querySelector(".documents-license__license-duration");if(s[t].enabled){if(n&&n.classList.contains("documents-license__license-icon--disabled")){n.classList.remove("documents-license__license-icon--disabled")}if(a&&a.classList.contains("documents-license__license-name--disabled")){a.classList.remove("documents-license__license-name--disabled")}if(e){e.textContent="Действует до: "+s[t].date;if(e.classList.contains("documents-license__license-duration--disabled")){e.classList.remove("documents-license__license-duration--disabled")}}}else{if(n&&!n.classList.contains("documents-license__license-icon--disabled")){n.classList.add("documents-license__license-icon--disabled")}if(a&&!a.classList.contains("documents-license__license-name--disabled")){a.classList.add("documents-license__license-name--disabled")}if(e){e.textContent="Отсутствует";if(!e.classList.contains("documents-license__license-duration--disabled")){e.classList.add("documents-license__license-duration--disabled")}}}});]]
    return code
end

function evalanon(code) evalcef(("(() => {%s})()"):format(code)) end

function evalcef(code, encoded)
    encoded = encoded or 0
    local bs = raknetNewBitStream()
    raknetBitStreamWriteInt8(bs, 17)
    raknetBitStreamWriteInt32(bs, 0)
    raknetBitStreamWriteInt16(bs, #code)
    raknetBitStreamWriteInt8(bs, encoded)
    raknetBitStreamWriteString(bs, code)
    raknetEmulPacketReceiveBitStream(220, bs)
    raknetDeleteBitStream(bs)
end

