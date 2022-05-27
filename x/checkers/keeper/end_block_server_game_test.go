package keeper_test

import (
	"testing"
	"time"

	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestForfeitUnplayed(t *testing.T) {
	_, keeper, context := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	game1, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	game1.Deadline = types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	keeper.SetStoredGame(ctx, game1)
	keeper.ForfeitExpiredGames(context)

	_, found = keeper.GetStoredGame(ctx, "1")
	require.False(t, found)

	nextGame, found := keeper.GetNextGame(ctx)
	require.True(t, found)
	require.EqualValues(t, types.NextGame{
		Creator:  "",
		IdValue:  2,
		FifoHead: "-1",
		FifoTail: "-1",
	}, nextGame)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 1)
	event := events[0]
	require.Equal(t, event.Type, "message")
	require.EqualValues(t, []sdk.Attribute{
		{Key: "module", Value: "checkers"},
		{Key: "action", Value: "GameForfeited"},
		{Key: "IdValue", Value: "1"},
		{Key: "Winner", Value: "*"},
	}, event.Attributes[6:])
}

func TestForfeitOlderUnplayed(t *testing.T) {
	msgServer, keeper, context := setupMsgServerWithOneGameForPlayMove(t)
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: bob,
		Red:     carol,
		Black:   alice,
	})
	ctx := sdk.UnwrapSDKContext(context)
	game1, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	game1.Deadline = types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	keeper.SetStoredGame(ctx, game1)
	keeper.ForfeitExpiredGames(context)

	_, found = keeper.GetStoredGame(ctx, "1")
	require.False(t, found)

	nextGame, found := keeper.GetNextGame(ctx)
	require.True(t, found)
	require.EqualValues(t, types.NextGame{
		Creator:  "",
		IdValue:  3,
		FifoHead: "2",
		FifoTail: "2",
	}, nextGame)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 1)
	event := events[0]
	require.Equal(t, event.Type, "message")
	require.EqualValues(t, []sdk.Attribute{
		{Key: "module", Value: "checkers"},
		{Key: "action", Value: "GameForfeited"},
		{Key: "IdValue", Value: "1"},
		{Key: "Winner", Value: "*"},
	}, event.Attributes[12:])
}

func TestForfeit2OldestUnplayedIn1Call(t *testing.T) {
	msgServer, keeper, context := setupMsgServerWithOneGameForPlayMove(t)
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: bob,
		Red:     carol,
		Black:   alice,
	})
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: carol,
		Red:     alice,
		Black:   bob,
	})
	ctx := sdk.UnwrapSDKContext(context)
	game1, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	game1.Deadline = types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	keeper.SetStoredGame(ctx, game1)
	game2, found := keeper.GetStoredGame(ctx, "2")
	require.True(t, found)
	game2.Deadline = types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	keeper.SetStoredGame(ctx, game2)
	keeper.ForfeitExpiredGames(context)

	_, found = keeper.GetStoredGame(ctx, "1")
	require.False(t, found)
	_, found = keeper.GetStoredGame(ctx, "2")
	require.False(t, found)

	nextGame, found := keeper.GetNextGame(ctx)
	require.True(t, found)
	require.EqualValues(t, types.NextGame{
		Creator:  "",
		IdValue:  4,
		FifoHead: "3",
		FifoTail: "3",
	}, nextGame)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 1)
	event := events[0]
	require.Equal(t, event.Type, "message")
	require.EqualValues(t, []sdk.Attribute{
		{Key: "module", Value: "checkers"},
		{Key: "action", Value: "GameForfeited"},
		{Key: "IdValue", Value: "1"},
		{Key: "Winner", Value: "*"},
	}, event.Attributes[18:22])
	require.EqualValues(t, []sdk.Attribute{
		{Key: "module", Value: "checkers"},
		{Key: "action", Value: "GameForfeited"},
		{Key: "IdValue", Value: "2"},
		{Key: "Winner", Value: "*"},
	}, event.Attributes[22:])
}

func TestForfeitPlayedOnce(t *testing.T) {
	msgServer, keeper, context := setupMsgServerWithOneGameForPlayMove(t)
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator: carol,
		IdValue: "1",
		FromX:   1,
		FromY:   2,
		ToX:     2,
		ToY:     3,
	})
	ctx := sdk.UnwrapSDKContext(context)
	game1, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	game1.Deadline = types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	keeper.SetStoredGame(ctx, game1)
	keeper.ForfeitExpiredGames(context)

	_, found = keeper.GetStoredGame(ctx, "1")
	require.False(t, found)

	nextGame, found := keeper.GetNextGame(ctx)
	require.True(t, found)
	require.EqualValues(t, types.NextGame{
		Creator:  "",
		IdValue:  2,
		FifoHead: "-1",
		FifoTail: "-1",
	}, nextGame)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 1)
	event := events[0]
	require.Equal(t, event.Type, "message")
	require.EqualValues(t, []sdk.Attribute{
		{Key: "module", Value: "checkers"},
		{Key: "action", Value: "GameForfeited"},
		{Key: "IdValue", Value: "1"},
		{Key: "Winner", Value: "*"},
	}, event.Attributes[13:])
}

func TestForfeitOlderPlayedOnce(t *testing.T) {
	msgServer, keeper, context := setupMsgServerWithOneGameForPlayMove(t)
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator: carol,
		IdValue: "1",
		FromX:   1,
		FromY:   2,
		ToX:     2,
		ToY:     3,
	})
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: bob,
		Red:     carol,
		Black:   alice,
	})
	ctx := sdk.UnwrapSDKContext(context)
	game1, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	game1.Deadline = types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	keeper.SetStoredGame(ctx, game1)
	keeper.ForfeitExpiredGames(context)

	_, found = keeper.GetStoredGame(ctx, "1")
	require.False(t, found)

	nextGame, found := keeper.GetNextGame(ctx)
	require.True(t, found)
	require.EqualValues(t, types.NextGame{
		Creator:  "",
		IdValue:  3,
		FifoHead: "2",
		FifoTail: "2",
	}, nextGame)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 1)
	event := events[0]
	require.Equal(t, event.Type, "message")
	require.EqualValues(t, []sdk.Attribute{
		{Key: "module", Value: "checkers"},
		{Key: "action", Value: "GameForfeited"},
		{Key: "IdValue", Value: "1"},
		{Key: "Winner", Value: "*"},
	}, event.Attributes[19:])
}

func TestForfeit2OldestPlayedOnceIn1Call(t *testing.T) {
	msgServer, keeper, context := setupMsgServerWithOneGameForPlayMove(t)
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator: carol,
		IdValue: "1",
		FromX:   1,
		FromY:   2,
		ToX:     2,
		ToY:     3,
	})
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: bob,
		Red:     carol,
		Black:   alice,
	})
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator: alice,
		IdValue: "2",
		FromX:   1,
		FromY:   2,
		ToX:     2,
		ToY:     3,
	})
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: carol,
		Red:     alice,
		Black:   bob,
	})
	ctx := sdk.UnwrapSDKContext(context)
	game1, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	game1.Deadline = types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	keeper.SetStoredGame(ctx, game1)
	game2, found := keeper.GetStoredGame(ctx, "2")
	require.True(t, found)
	game2.Deadline = types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	keeper.SetStoredGame(ctx, game2)
	keeper.ForfeitExpiredGames(context)

	_, found = keeper.GetStoredGame(ctx, "1")
	require.False(t, found)
	_, found = keeper.GetStoredGame(ctx, "2")
	require.False(t, found)

	nextGame, found := keeper.GetNextGame(ctx)
	require.True(t, found)
	require.EqualValues(t, types.NextGame{
		Creator:  "",
		IdValue:  4,
		FifoHead: "3",
		FifoTail: "3",
	}, nextGame)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 1)
	event := events[0]
	require.Equal(t, event.Type, "message")
	require.EqualValues(t, []sdk.Attribute{
		{Key: "module", Value: "checkers"},
		{Key: "action", Value: "GameForfeited"},
		{Key: "IdValue", Value: "1"},
		{Key: "Winner", Value: "*"},
	}, event.Attributes[32:36])
	require.EqualValues(t, []sdk.Attribute{
		{Key: "module", Value: "checkers"},
		{Key: "action", Value: "GameForfeited"},
		{Key: "IdValue", Value: "2"},
		{Key: "Winner", Value: "*"},
	}, event.Attributes[36:])
}

func TestForfeitPlayedTwice(t *testing.T) {
	msgServer, keeper, context := setupMsgServerWithOneGameForPlayMove(t)
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator: carol,
		IdValue: "1",
		FromX:   1,
		FromY:   2,
		ToX:     2,
		ToY:     3,
	})
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator: bob,
		IdValue: "1",
		FromX:   0,
		FromY:   5,
		ToX:     1,
		ToY:     4,
	})
	ctx := sdk.UnwrapSDKContext(context)
	game1, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	oldDeadline := types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	game1.Deadline = oldDeadline
	keeper.SetStoredGame(ctx, game1)
	keeper.ForfeitExpiredGames(context)

	game1, found = keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Creator:   alice,
		Index:     "1",
		Game:      "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|*r******|**r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:      "b",
		Red:       bob,
		Black:     carol,
		MoveCount: uint64(2),
		BeforeId:  "-1",
		AfterId:   "-1",
		Deadline:  oldDeadline,
		Winner:    "r",
	}, game1)

	nextGame, found := keeper.GetNextGame(ctx)
	require.True(t, found)
	require.EqualValues(t, types.NextGame{
		Creator:  "",
		IdValue:  2,
		FifoHead: "-1",
		FifoTail: "-1",
	}, nextGame)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 1)
	event := events[0]
	require.Equal(t, event.Type, "message")
	require.EqualValues(t, []sdk.Attribute{
		{Key: "module", Value: "checkers"},
		{Key: "action", Value: "GameForfeited"},
		{Key: "IdValue", Value: "1"},
		{Key: "Winner", Value: "r"},
	}, event.Attributes[20:])
}

func TestForfeitOlderPlayedTwice(t *testing.T) {
	msgServer, keeper, context := setupMsgServerWithOneGameForPlayMove(t)
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator: carol,
		IdValue: "1",
		FromX:   1,
		FromY:   2,
		ToX:     2,
		ToY:     3,
	})
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator: bob,
		IdValue: "1",
		FromX:   0,
		FromY:   5,
		ToX:     1,
		ToY:     4,
	})
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: bob,
		Red:     carol,
		Black:   alice,
	})
	ctx := sdk.UnwrapSDKContext(context)
	game1, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	oldDeadline := types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	game1.Deadline = oldDeadline
	keeper.SetStoredGame(ctx, game1)
	keeper.ForfeitExpiredGames(context)

	game1, found = keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Creator:   alice,
		Index:     "1",
		Game:      "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|*r******|**r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:      "b",
		Red:       bob,
		Black:     carol,
		MoveCount: uint64(2),
		BeforeId:  "-1",
		AfterId:   "-1",
		Deadline:  oldDeadline,
		Winner:    "r",
	}, game1)

	nextGame, found := keeper.GetNextGame(ctx)
	require.True(t, found)
	require.EqualValues(t, types.NextGame{
		Creator:  "",
		IdValue:  3,
		FifoHead: "2",
		FifoTail: "2",
	}, nextGame)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 1)
	event := events[0]
	require.Equal(t, event.Type, "message")
	require.EqualValues(t, []sdk.Attribute{
		{Key: "module", Value: "checkers"},
		{Key: "action", Value: "GameForfeited"},
		{Key: "IdValue", Value: "1"},
		{Key: "Winner", Value: "r"},
	}, event.Attributes[26:])
}

func TestForfeit2OldestPlayedTwiceIn1Call(t *testing.T) {
	msgServer, keeper, context := setupMsgServerWithOneGameForPlayMove(t)
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator: carol,
		IdValue: "1",
		FromX:   1,
		FromY:   2,
		ToX:     2,
		ToY:     3,
	})
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator: bob,
		IdValue: "1",
		FromX:   0,
		FromY:   5,
		ToX:     1,
		ToY:     4,
	})
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: bob,
		Red:     carol,
		Black:   alice,
	})
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator: alice,
		IdValue: "2",
		FromX:   1,
		FromY:   2,
		ToX:     2,
		ToY:     3,
	})
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator: carol,
		IdValue: "2",
		FromX:   0,
		FromY:   5,
		ToX:     1,
		ToY:     4,
	})
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: carol,
		Red:     alice,
		Black:   bob,
	})
	ctx := sdk.UnwrapSDKContext(context)
	game1, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	oldDeadline := types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	game1.Deadline = oldDeadline
	keeper.SetStoredGame(ctx, game1)
	game2, found := keeper.GetStoredGame(ctx, "2")
	require.True(t, found)
	game2.Deadline = oldDeadline
	keeper.SetStoredGame(ctx, game2)
	keeper.ForfeitExpiredGames(context)

	game1, found = keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Creator:   alice,
		Index:     "1",
		Game:      "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|*r******|**r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:      "b",
		Red:       bob,
		Black:     carol,
		MoveCount: uint64(2),
		BeforeId:  "-1",
		AfterId:   "-1",
		Deadline:  oldDeadline,
		Winner:    "r",
	}, game1)

	game2, found = keeper.GetStoredGame(ctx, "2")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Creator:   bob,
		Index:     "2",
		Game:      "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|*r******|**r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:      "b",
		Red:       carol,
		Black:     alice,
		MoveCount: uint64(2),
		BeforeId:  "-1",
		AfterId:   "-1",
		Deadline:  oldDeadline,
		Winner:    "r",
	}, game2)

	nextGame, found := keeper.GetNextGame(ctx)
	require.True(t, found)
	require.EqualValues(t, types.NextGame{
		Creator:  "",
		IdValue:  4,
		FifoHead: "3",
		FifoTail: "3",
	}, nextGame)

	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 1)
	event := events[0]
	require.Equal(t, event.Type, "message")
	require.EqualValues(t, []sdk.Attribute{
		{Key: "module", Value: "checkers"},
		{Key: "action", Value: "GameForfeited"},
		{Key: "IdValue", Value: "1"},
		{Key: "Winner", Value: "r"},
	}, event.Attributes[46:50])
	require.EqualValues(t, []sdk.Attribute{
		{Key: "module", Value: "checkers"},
		{Key: "action", Value: "GameForfeited"},
		{Key: "IdValue", Value: "2"},
		{Key: "Winner", Value: "r"},
	}, event.Attributes[50:])
}
