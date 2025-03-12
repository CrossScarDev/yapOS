import("CoreLibs/object")
import("CoreLibs/graphics")
import("CoreLibs/sprites")
import("CoreLibs/timer")

local gfx = playdate.graphics

local games = {}
for _, group in ipairs(playdate.datastore.read("/System/Data/games")) do
	for _, game in ipairs(group["games"]) do
		table.insert(games, game)
	end
end

local selectedGame = 1

local font = gfx.font.new("fonts/roobert11")
gfx.setFont(font)

function playdate.update()
	gfx.clear()

	local y = 20
	for i, game in ipairs(games) do
		gfx.setImageDrawMode(playdate.graphics.kDrawModeFillBlack)
		if i == selectedGame then
			gfx.setColor(playdate.graphics.kColorBlack)
			gfx.fillRoundRect(10, y - 5, font:getTextWidth(game["title"]) + 20, font:getHeight() + 10, 5)
			gfx.setImageDrawMode(playdate.graphics.kDrawModeFillWhite)
		end
		gfx.drawText(game["title"], 20, y)
		y += font:getHeight() + 20
	end
end

function playdate.downButtonDown()
	if selectedGame < #games then
		selectedGame += 1
	end
end

function playdate.upButtonDown()
	if selectedGame > 1 then
		selectedGame -= 1
	end
end

function playdate.AButtonDown()
	playdate.system.switchToGame(games[selectedGame]["path"])
end
