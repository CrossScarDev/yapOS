import("CoreLibs/object")
import("CoreLibs/graphics")
import("CoreLibs/sprites")
import("CoreLibs/timer")
import("CoreLibs/ui")

local gfx <const> = playdate.graphics
local ui <const> = playdate.ui

local games = {}
for _, group in ipairs(playdate.datastore.read("/System/Data/games")) do
	for _, game in ipairs(group["games"]) do
		table.insert(games, game)
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
		games[row]["title"],
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
	playdate.system.switchToGame(games[listview:getSelectedRow()]["path"])
end
