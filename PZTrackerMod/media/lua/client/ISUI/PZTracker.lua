-- #########################
-- # PZTrackerMod by BigJk #
-- #########################

local interval = 10 -- Set the interval here! Small number like 10 is really fast. Higher numbers get slower.

local currentInterval = 0;
local enableFollow = false;

-- CALLBACK'S --
function onGameStart()
	-- Print out some basic player infos on game start. Maybe used in later and more advanced version.
	local player = getSpecificPlayer(0);
	emitEvent("BasicInfo", {
		Username = player:getUsername(),
		Surname = player:getSurname(),
		Name = player:getForname()
	});
end

function onConnected()
	emitBasicEvent("Connected")
end

function onDisconnect()
	emitBasicEvent("Disconnect")
end

function onMove()
	-- Check if follow-mode is enabled.
	if enableFollow == false then
		return;
	end

	-- Only performe update every X calls to reduce performance cost. Interval can be set above.
	currentInterval = currentInterval + 1;
	if currentInterval ~= interval then
		return;
	end
	currentInterval = 0;
	
	-- Get local player and pring X and Y position.
    local player = getSpecificPlayer(0);
	if player then
		emitEvent("Position", {
			X = player:getX(),
			Y = player:getY()
		});
	end
end

function onkeyPressed(num)
	if num == 200 then -- UP : Move map up.
		emitBasicEvent("UpPress");
	elseif num == 203 then -- LEFT : Move map left.
		emitBasicEvent("LeftPress");
	elseif num == 205 then -- RIGHT : Move map right.
		emitBasicEvent("RightPress");
	elseif num == 208 then -- DOWN : Move map down.
		emitBasicEvent("DownPress");
	elseif num == 33 then -- F : Enable follow-mode.
		enableFollow = not enableFollow;
	elseif num == 78 then -- + : Zoom in.
		emitBasicEvent("ZoomInPress");
	elseif num == 74 then -- - : Zoom out.
		emitBasicEvent("ZoomOutPress");
	end
end

-- BASIC FUNCTION'S --
function emitBasicEvent(t) -- Emits a simple event.
	print("Tracker:" .. t);
end

function emitEvent(name, data) -- Emits a event that contains data.
	eventStr = "Tracker:" .. name;
	for k, v in pairs(data) do
		eventStr = eventStr .. " " .. k .. "(" .. v .. ")";
	end
	print(eventStr);
end


-- EVENT'S --
Events.OnGameStart.Add(onGameStart);
Events.OnPlayerMove.Add(onMove);
Events.OnKeyPressed.Add(onkeyPressed);
Events.OnConnected.Add(onConnected);
Events.OnDisconnect.Add(onDisconnect);