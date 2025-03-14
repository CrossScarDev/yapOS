import("CoreLibs/object")
import("CoreLibs/graphics")
import("CoreLibs/sprites")
import("CoreLibs/timer")
import("CoreLibs/ui")

local gfx <const> = playdate.graphics
local ui <const> = playdate.ui
local sys <const> = playdate.system

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
gfx.setFont(font)

local listview = ui.gridview.new(playdate.display.getWidth() / 2, font:getHeight() + 20)
listview:setNumberOfRows(#games)
listview:setNumberOfColumns(1)
listview:setCellPadding(0, 0, 2.5, 2.5)
listview:setContentInset(5, 0, 5, 0)
listview:setSelectedRow(1)
function listview:drawCell(section, row, column, selected, x, y, width, height)
	gfx.setImageDrawMode(playdate.graphics.kDrawModeFillBlack)
	if selected then
		gfx.setColor(gfx.kColorBlack)
		gfx.fillRoundRect(x, y, width, height, height / 2)
		gfx.setImageDrawMode(playdate.graphics.kDrawModeFillWhite)
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
	gfx.clear()
	playdate.timer.updateTimers()
	listview:drawInRect(0, 0, playdate.display.getWidth(), playdate.display.getHeight())

	local selectedGame = games[listview:getSelectedRow()]
	local infoOffset = 5
	local sideOffset = playdate.display.getWidth() / 2 + 10
	if selectedGame.path and selectedGame.imagepath then
		if playdate.file.exists(selectedGame.path .. "/" .. selectedGame.imagepath .. "/icon.pdi") then
			local icon = gfx.image.new(selectedGame.path .. "/" .. selectedGame.imagepath .. "/icon.pdi")
			icon:drawScaled(sideOffset, 5, 2)
			local w, h = gfx.getTextSizeForMaxWidth(selectedGame.name, playdate.display.getWidth() - sideOffset + 59)
			gfx.drawTextInRect(
				selectedGame.name,
				sideOffset + 74,
				37 - h / 2,
				playdate.display.getWidth() - sideOffset + 59,
				64
			)

			infoOffset += 69
		end
	end

	if selectedGame.version then
		gfx.drawText("Version: " .. selectedGame.version, sideOffset, infoOffset)
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
