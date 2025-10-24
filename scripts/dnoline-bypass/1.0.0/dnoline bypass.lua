-- CTEA by SanTrope RP (modified https://github.com/somesocks/lua-lockbox/blob/master/lockbox/cipher/tea.lua)
local bit = require('bit')
local AND = bit.band;
local OR  = bit.bor;
local XOR = bit.bxor;
local LSHIFT = bit.lshift;
local RSHIFT = bit.rshift;

local TEA = {};
TEA.blockSize = 8;

TEA.encrypt = function(key, data)
    local y = data[1];
    local z = data[2];
    local delta = 0x9e3779b9;
    local sum = 0;

    local k0 = key[1];
    local k1 = key[2];
    local k2 = key[3];
    local k3 = key[4];

    for _ = 1, 32 do
        local temp;

        sum = AND(sum + delta, 0xFFFFFFFF);

        temp = z + sum;
        temp = XOR(temp, LSHIFT(z, 4) + k0);
        temp = XOR(temp, RSHIFT(z, 5) + k1);
        y = AND(y + temp, 0xFFFFFFFF);

        temp = y + sum;
        temp = XOR(temp, LSHIFT(y, 4) + k2);
        temp = XOR(temp, RSHIFT(y, 5) + k3);
        z = AND( z + temp, 0xFFFFFFFF);
    end

    local out = {};

    out[1] = y;
    out[2] = z;

    return out;
end

-- ��� �����
local sampev = require 'samp.events'

local gpci = 'FF2BE5E6F5D9392F57C4E66F7AD78767277C6E4F6B'
local version = "0.3.7"

local onlineCamera = 1 -- island
local onlineConstant = 6
local onlineVersion = 509

local online = false
local sentMessageStatsSync = false

function main()
    repeat wait(0) until isSampAvailable()

    local ip, port = sampGetCurrentServerAddress()
    if      string.find(ip, "194.147.32.192") == nil
        and string.find(ip, "194.147.32.193") == nil
        and string.find(ip, "80.66.71.65")    == nil
        and string.find(ip, "80.66.71.64")    == nil
        and string.find(ip, "80.66.71.63")    == nil
    then
        printStringNow('~r~online bypass disabled', 2000)
    else
        online = true
        sampAddChatMessage("����� ��� �������� ���������� ������� 666.228", -1)
    end

    wait(-1)
end

function sampev.onSendClientJoin(v, m, nickname, challengeResponse, a, c, u)
    if online then
        -- ������� �������� �������� �����: ���� ���������
        for i = 1, 5 do
            local bs = raknetNewBitStream()
            raknetBitStreamWriteInt32(bs, onlineVersion)
            raknetBitStreamWriteInt8(bs, onlineCamera)
            raknetSendRpcEx(140, bs, 1, 10, 0, false)
            raknetDeleteBitStream(bs)
        end

        local bs = raknetNewBitStream()
        raknetBitStreamWriteInt32(bs, 0xFD9)
        raknetBitStreamWriteInt8(bs, 1)
        raknetBitStreamWriteInt8(bs, #nickname)
        raknetBitStreamWriteString(bs, nickname)
        raknetBitStreamWriteInt32(bs, challengeResponse)
        raknetBitStreamWriteInt8(bs, #gpci)
        raknetBitStreamWriteString(bs, gpci)
        raknetBitStreamWriteInt8(bs, #version)
        raknetBitStreamWriteString(bs, version)
        raknetBitStreamWriteInt8(bs, onlineConstant)
        raknetBitStreamWriteInt32(bs, onlineVersion)
        raknetSendRpcEx(25, bs, 1, 10, 0, false)
        raknetDeleteBitStream(bs)

        sampAddChatMessage("����� ��� �������, ��� ��������...", -1)

        return false
    end
end

function onReceiveRpc(id, bitStream)
    if id == 87 and online then
        local textdrawId = raknetBitStreamReadInt16(bitStream)

        if textdrawId == 0xFFFF then
            local keys = {0, 0, 0, 0}
            local data = {0, 0}

            for i = 1, 4 do
                keys[i] = raknetBitStreamReadInt32(bitStream)
            end

            for i = 1, 2 do
                data[i] = raknetBitStreamReadInt32(bitStream)
            end

            local pomoika = TEA.encrypt(keys, data)

            local bs = raknetNewBitStream()
            raknetBitStreamWriteInt16(bs, 0xFFFF)

            for i = 1, 17 do
                raknetBitStreamWriteInt8(bs, math.random(1, 255))
            end

            raknetBitStreamWriteInt32(bs, pomoika[1])
            raknetBitStreamWriteInt32(bs, pomoika[2])
            raknetSendRpcEx(83, bs, 0, 9, 0, false)
            raknetDeleteBitStream(bs)

            sampAddChatMessage("����� online ����� ������ ��������, ����� mhertz", -1)
            return false
        end
    end
end

function onSendPacket(id, bs, priority, reliability, orderingChannel)
    if id == 205 and online then
        raknetBitStreamSetWriteOffset(bs, raknetBitStreamGetNumberOfBitsUsed(bs))
        raknetBitStreamWriteInt8(bs, 100) -- breath level (0 - 100)

        if not sentMessageStatsSync then
            sampAddChatMessage("�� ����� ����� ��� ������ ������� �������� ��� �� �� ������ � ��� �� ���� � ���� ���", -1)
            sentMessageStatsSync = true
        end
    end
end