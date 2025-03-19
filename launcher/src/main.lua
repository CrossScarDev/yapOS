import("CoreLibs/object")
import("CoreLibs/graphics")
import("CoreLibs/sprites")
import("CoreLibs/timer")
import("CoreLibs/ui")

local gfx <const> = playdate.graphics
local ui <const> = playdate.ui
local sys <const> = playdate.system
local menu <const> = playdate.getSystemMenu()

local games = {}

local scrollSoundUp = playdate.sound.fileplayer.new("systemsfx/01-selection-trimmed")
local scrollSoundDown = playdate.sound.fileplayer.new("systemsfx/02-selection-reverse-trimmed")

for _, group in ipairs(sys.getInstalledGameList()) do
	for _, game in ipairs(group) do
		if sys.game.getPath(game) then
			local path = sys.game.getPath(game)
			if sys.game.getBundleID(game) then
				local props = sys.getMetadata(path .. "/pdxinfo")
				local newprops = {}
				for k, v in pairs(props) do
					newprops[string.lower(k)] = v
				end
				props = newprops
				props["path"] = path
				props["group"] = group.name
				props["suppressContentWarning"] = game:getSuppressContentWarning()
				table.insert(games, props)
			end
		end
	end
end

local font = gfx.font.new("fonts/roobert11")
local fontBold = gfx.font.new("fonts/roobert11Bold")
gfx.setFont(font)
gfx.setFont(fontBold, gfx.font.kVariantBold)

local settings = playdate.datastore.read("/System/Data/yapOS")
local darkMode = false
local invertIcons = false

if settings == nil or settings.dark == nil then
	playdate.datastore.write({ ["dark"] = darkMode, ["invert"] = invertIcons }, "/System/Data/yapOS")
end

if settings ~= nil then
	if settings.dark ~= nil then
		darkMode = settings.dark
	end
	if settings.invert ~= nil then
		invertIcons = settings.invert
	end
end

local listview = ui.gridview.new(playdate.display.getWidth() / 2 - 7.5, font:getHeight() + 20)
listview:setNumberOfRows(#games)
listview:setNumberOfColumns(1)
listview:setCellPadding(0, 0, 2.5, 2.5)
listview:setContentInset(5, 0, 5, 0)
listview:setSelectedRow(1)
function listview:drawCell(_, row, _, selected, x, y, width, height)
	gfx.setImageDrawMode(darkMode and gfx.kDrawModeFillWhite or gfx.kDrawModeFillBlack)
	if selected then
		gfx.setColor(darkMode and gfx.kColorWhite or gfx.kColorBlack)
		gfx.fillRoundRect(x, y, width, height, height / 2)
		gfx.setImageDrawMode(darkMode and gfx.kDrawModeFillBlack or gfx.kDrawModeFillWhite)
	end
	gfx.drawTextInRect(
		games[row].name,
		x + 10,
		y + height / 2 - font:getHeight() / 2,
		width - 20,
		height,
		nil,
		nil,
		kTextAlignment.center
	)
end

local function renderLeft()
	gfx.setColor(darkMode and gfx.kColorBlack or gfx.kColorWhite)
	gfx.fillRect(0, 0, playdate.display.getWidth() / 2, playdate.display.getHeight())
	listview:drawInRect(0, 0, playdate.display.getWidth() / 2, playdate.display.getHeight())

	playdate.setMenuImage(gfx.getDisplayImage(), playdate.display.getWidth() / 2 - 2)
end

local function renderRight()
	gfx.setColor(darkMode and gfx.kColorBlack or gfx.kColorWhite)
	gfx.fillRect(playdate.display.getWidth() / 2, 0, playdate.display.getWidth() / 2, playdate.display.getHeight())

	local selectedGame = games[listview:getSelectedRow()]
	local infoOffset = 5
	local sideOffset = playdate.display.getWidth() / 2 + 2.5
	if
		(
			selectedGame.path
			and selectedGame.imagepath
			and playdate.file.exists(selectedGame.path .. "/" .. selectedGame.imagepath .. "/icon.pdi")
		)
		or (
			string.match(selectedGame.group, "^Season%-%d%d%d$")
			and playdate.file.exists("s1_icons/" .. selectedGame.bundleid .. ".pdi")
		)
	then
		gfx.setImageDrawMode(gfx.kDrawModeCopy)
		local icon = gfx.image.new(
			string.match(selectedGame.group, "^Season%-%d%d%d$") and ("s1_icons/" .. selectedGame.bundleid .. ".pdi")
				or (selectedGame.path .. "/" .. selectedGame.imagepath .. "/icon.pdi")
		)
		gfx.setColor(gfx.kColorWhite)
		gfx.fillRect(sideOffset, 5, 64, 64)
		gfx.setImageDrawMode(invertIcons and gfx.kDrawModeInverted or gfx.kDrawModeCopy)
		icon:drawScaled(sideOffset, 5, 2)
		gfx.setImageDrawMode(gfx.kDrawModeCopy)

		gfx.setImageDrawMode(darkMode and gfx.kDrawModeFillWhite or gfx.kDrawModeFillBlack)
		local w, h = gfx.getTextSizeForMaxWidth(selectedGame.name, playdate.display.getWidth() - sideOffset - 79)
		gfx.drawTextInRect(
			selectedGame.name,
			sideOffset + 74,
			h < 64 and 37 - h / 2 or 5,
			playdate.display.getWidth() - sideOffset - 79,
			64
		)

		infoOffset += 69
	else
		infoOffset += 5
		gfx.setImageDrawMode(darkMode and gfx.kDrawModeFillWhite or gfx.kDrawModeFillBlack)
	end

	gfx.drawTextInRect(
		"*Author:* " .. selectedGame.author,
		sideOffset,
		infoOffset,
		playdate.display.getWidth() - sideOffset - 5,
		font:getHeight(),
		nil,
		"..."
	)
	infoOffset += font:getHeight() + 5
	if selectedGame.version then
		gfx.drawTextInRect(
			"*Version:* " .. selectedGame.version,
			sideOffset,
			infoOffset,
			playdate.display.getWidth() - sideOffset - 5,
			font:getHeight(),
			nil,
			"..."
		)
		infoOffset += font:getHeight() + 5
	end
	gfx.drawTextInRect(
		"*Group:* " .. selectedGame.group,
		sideOffset,
		infoOffset,
		playdate.display.getWidth() - sideOffset - 5,
		font:getHeight(),
		nil,
		"..."
	)
	infoOffset += font:getHeight() + 5
	if selectedGame.description then
		gfx.drawTextInRect(
			"*Descripton:*\n" .. selectedGame.description,
			sideOffset,
			infoOffset,
			playdate.display.getWidth() - sideOffset - 5,
			playdate.display.getHeight() - infoOffset - 5,
			nil,
			"..."
		)
	end
end

local frameCount = 0

function playdate.update()
	playdate.timer.updateTimers()
	if listview.needsDisplay then
		renderLeft()
	end
	if frameCount < 1 then
		frameCount += 1
	elseif frameCount == 1 then
		playdate.setMenuImage(gfx.getDisplayImage(), playdate.display.getWidth() / 2 - 2)
	end
end

function playdate.downButtonDown()
	if listview:getSelectedRow() < #games then
		listview:setSelectedRow(listview:getSelectedRow() + 1)
		listview:scrollCellToCenter(1, listview:getSelectedRow(), 1)
		scrollSoundDown:play()
		renderRight()
	end
end

function playdate.upButtonDown()
	if listview:getSelectedRow() > 1 then
		listview:setSelectedRow(listview:getSelectedRow() - 1)
		listview:scrollCellToCenter(1, listview:getSelectedRow(), 1)
		scrollSoundUp:play()
		renderRight()
	end
end

function playdate.AButtonDown()
	sys.switchToGame(games[listview:getSelectedRow()]["path"])
end

menu:addCheckmarkMenuItem("dark mode", darkMode, function(v)
	darkMode = v
	playdate.datastore.write({ ["dark"] = darkMode }, "/System/Data/yapOS")
	renderLeft()
	renderRight()
end)
menu:addCheckmarkMenuItem("invert icons", invertIcons, function(v)
	invertIcons = v
	playdate.datastore.write({ ["invert"] = invertIcons }, "/System/Data/yapOS")
	renderRight()
end)

renderLeft()
renderRight()
