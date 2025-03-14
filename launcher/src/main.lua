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

local darkMode = false

menu:addCheckmarkMenuItem("Dark Mode", darkMode, function(v)
	darkMode = v
end)

local listview = ui.gridview.new(playdate.display.getWidth() / 2, font:getHeight() + 20)
listview:setNumberOfRows(#games)
listview:setNumberOfColumns(1)
listview:setCellPadding(0, 0, 2.5, 2.5)
listview:setContentInset(5, 0, 5, 0)
listview:setSelectedRow(1)
function listview:drawCell(section, row, column, selected, x, y, width, height)
	gfx.setImageDrawMode(darkMode and gfx.kDrawModeFillWhite or gfx.kDrawModeFillBlack)
	if selected then
		gfx.setColor(darkMode and gfx.kColorWhite or gfx.kColorBlack)
		gfx.fillRoundRect(x, y, width, height, height / 2)
		gfx.setImageDrawMode(darkMode and gfx.kDrawModeFillBlack or gfx.kDrawModeFillWhite)
	end
	gfx.drawTextInRect(
		games[row].name,
		x,
		y + height / 2 - font:getHeight() / 2,
		width,
		height,
		nil,
		nil,
		kTextAlignment.center
	)
end

function playdate.update()
	gfx.clear(darkMode and gfx.kColorBlack or gfx.kColorWhite)
	playdate.timer.updateTimers()
	listview:drawInRect(0, 0, playdate.display.getWidth(), playdate.display.getHeight())

	local selectedGame = games[listview:getSelectedRow()]
	local infoOffset = 5
	local sideOffset = playdate.display.getWidth() / 2 + 10
	if
		selectedGame.path
		and selectedGame.imagepath
		and playdate.file.exists(selectedGame.path .. "/" .. selectedGame.imagepath .. "/icon.pdi")
	then
		gfx.setImageDrawMode(gfx.kDrawModeCopy)
		local icon = gfx.image.new(selectedGame.path .. "/" .. selectedGame.imagepath .. "/icon.pdi")
		gfx.setColor(gfx.kColorWhite)
		gfx.fillRect(sideOffset, 5, 64, 64)
		icon:drawScaled(sideOffset, 5, 2)

		gfx.setImageDrawMode(darkMode and gfx.kDrawModeFillWhite or gfx.kDrawModeFillBlack)
		local w, h = gfx.getTextSizeForMaxWidth(selectedGame.name, playdate.display.getWidth() - sideOffset + 59)
		gfx.drawTextInRect(
			selectedGame.name,
			sideOffset + 74,
			h < 64 and 37 - h / 2 or 5,
			playdate.display.getWidth() - sideOffset + 59,
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

function playdate.downButtonDown()
	if listview:getSelectedRow() < #games then
		listview:setSelectedRow(listview:getSelectedRow() + 1)
		listview:scrollCellToCenter(1, listview:getSelectedRow(), 1)
	end
end

function playdate.upButtonDown()
	if listview:getSelectedRow() > 1 then
		listview:setSelectedRow(listview:getSelectedRow() - 1)
		listview:scrollCellToCenter(1, listview:getSelectedRow(), 1)
	end
end

function playdate.AButtonDown()
	sys.switchToGame(games[listview:getSelectedRow()]["path"])
end
