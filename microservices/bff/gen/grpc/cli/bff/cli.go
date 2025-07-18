// Code generated by goa v3.21.1, DO NOT EDIT.
//
// bff gRPC client CLI support package
//
// Command:
// $ goa gen object-t.com/hackz-giganoto/microservices/bff/design

package cli

import (
	"flag"
	"fmt"
	"os"

	goa "goa.design/goa/v3/pkg"
	grpc "google.golang.org/grpc"
	bffc "object-t.com/hackz-giganoto/microservices/bff/gen/grpc/bff/client"
)

// UsageCommands returns the set of commands and sub-commands using the format
//
//	command (subcommand1|subcommand2|...)
func UsageCommands() string {
	return `bff (create-room|history|room-list|join-room|invite-room|stream-chat|get-profile|update-profile)
`
}

// UsageExamples produces an example of a valid invocation of the CLI tool.
func UsageExamples() string {
	return os.Args[0] + ` bff create-room --token "Aspernatur tenetur libero accusantium laborum."` + "\n" +
		""
}

// ParseEndpoint returns the endpoint and payload as specified on the command
// line.
func ParseEndpoint(
	cc *grpc.ClientConn,
	opts ...grpc.CallOption,
) (goa.Endpoint, any, error) {
	var (
		bffFlags = flag.NewFlagSet("bff", flag.ContinueOnError)

		bffCreateRoomFlags     = flag.NewFlagSet("create-room", flag.ExitOnError)
		bffCreateRoomTokenFlag = bffCreateRoomFlags.String("token", "REQUIRED", "")

		bffHistoryFlags       = flag.NewFlagSet("history", flag.ExitOnError)
		bffHistoryMessageFlag = bffHistoryFlags.String("message", "", "")
		bffHistoryTokenFlag   = bffHistoryFlags.String("token", "REQUIRED", "")

		bffRoomListFlags     = flag.NewFlagSet("room-list", flag.ExitOnError)
		bffRoomListTokenFlag = bffRoomListFlags.String("token", "REQUIRED", "")

		bffJoinRoomFlags       = flag.NewFlagSet("join-room", flag.ExitOnError)
		bffJoinRoomMessageFlag = bffJoinRoomFlags.String("message", "", "")
		bffJoinRoomTokenFlag   = bffJoinRoomFlags.String("token", "REQUIRED", "")

		bffInviteRoomFlags       = flag.NewFlagSet("invite-room", flag.ExitOnError)
		bffInviteRoomMessageFlag = bffInviteRoomFlags.String("message", "", "")
		bffInviteRoomTokenFlag   = bffInviteRoomFlags.String("token", "REQUIRED", "")

		bffStreamChatFlags      = flag.NewFlagSet("stream-chat", flag.ExitOnError)
		bffStreamChatTokenFlag  = bffStreamChatFlags.String("token", "REQUIRED", "")
		bffStreamChatRoomIDFlag = bffStreamChatFlags.String("room-id", "REQUIRED", "")

		bffGetProfileFlags       = flag.NewFlagSet("get-profile", flag.ExitOnError)
		bffGetProfileMessageFlag = bffGetProfileFlags.String("message", "", "")
		bffGetProfileTokenFlag   = bffGetProfileFlags.String("token", "REQUIRED", "")

		bffUpdateProfileFlags       = flag.NewFlagSet("update-profile", flag.ExitOnError)
		bffUpdateProfileMessageFlag = bffUpdateProfileFlags.String("message", "", "")
		bffUpdateProfileTokenFlag   = bffUpdateProfileFlags.String("token", "", "")
	)
	bffFlags.Usage = bffUsage
	bffCreateRoomFlags.Usage = bffCreateRoomUsage
	bffHistoryFlags.Usage = bffHistoryUsage
	bffRoomListFlags.Usage = bffRoomListUsage
	bffJoinRoomFlags.Usage = bffJoinRoomUsage
	bffInviteRoomFlags.Usage = bffInviteRoomUsage
	bffStreamChatFlags.Usage = bffStreamChatUsage
	bffGetProfileFlags.Usage = bffGetProfileUsage
	bffUpdateProfileFlags.Usage = bffUpdateProfileUsage

	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		return nil, nil, err
	}

	if flag.NArg() < 2 { // two non flag args are required: SERVICE and ENDPOINT (aka COMMAND)
		return nil, nil, fmt.Errorf("not enough arguments")
	}

	var (
		svcn string
		svcf *flag.FlagSet
	)
	{
		svcn = flag.Arg(0)
		switch svcn {
		case "bff":
			svcf = bffFlags
		default:
			return nil, nil, fmt.Errorf("unknown service %q", svcn)
		}
	}
	if err := svcf.Parse(flag.Args()[1:]); err != nil {
		return nil, nil, err
	}

	var (
		epn string
		epf *flag.FlagSet
	)
	{
		epn = svcf.Arg(0)
		switch svcn {
		case "bff":
			switch epn {
			case "create-room":
				epf = bffCreateRoomFlags

			case "history":
				epf = bffHistoryFlags

			case "room-list":
				epf = bffRoomListFlags

			case "join-room":
				epf = bffJoinRoomFlags

			case "invite-room":
				epf = bffInviteRoomFlags

			case "stream-chat":
				epf = bffStreamChatFlags

			case "get-profile":
				epf = bffGetProfileFlags

			case "update-profile":
				epf = bffUpdateProfileFlags

			}

		}
	}
	if epf == nil {
		return nil, nil, fmt.Errorf("unknown %q endpoint %q", svcn, epn)
	}

	// Parse endpoint flags if any
	if svcf.NArg() > 1 {
		if err := epf.Parse(svcf.Args()[1:]); err != nil {
			return nil, nil, err
		}
	}

	var (
		data     any
		endpoint goa.Endpoint
		err      error
	)
	{
		switch svcn {
		case "bff":
			c := bffc.NewClient(cc, opts...)
			switch epn {
			case "create-room":
				endpoint = c.CreateRoom()
				data, err = bffc.BuildCreateRoomPayload(*bffCreateRoomTokenFlag)
			case "history":
				endpoint = c.History()
				data, err = bffc.BuildHistoryPayload(*bffHistoryMessageFlag, *bffHistoryTokenFlag)
			case "room-list":
				endpoint = c.RoomList()
				data, err = bffc.BuildRoomListPayload(*bffRoomListTokenFlag)
			case "join-room":
				endpoint = c.JoinRoom()
				data, err = bffc.BuildJoinRoomPayload(*bffJoinRoomMessageFlag, *bffJoinRoomTokenFlag)
			case "invite-room":
				endpoint = c.InviteRoom()
				data, err = bffc.BuildInviteRoomPayload(*bffInviteRoomMessageFlag, *bffInviteRoomTokenFlag)
			case "stream-chat":
				endpoint = c.StreamChat()
				data, err = bffc.BuildStreamChatPayload(*bffStreamChatTokenFlag, *bffStreamChatRoomIDFlag)
			case "get-profile":
				endpoint = c.GetProfile()
				data, err = bffc.BuildGetProfilePayload(*bffGetProfileMessageFlag, *bffGetProfileTokenFlag)
			case "update-profile":
				endpoint = c.UpdateProfile()
				data, err = bffc.BuildUpdateProfilePayload(*bffUpdateProfileMessageFlag, *bffUpdateProfileTokenFlag)
			}
		}
	}
	if err != nil {
		return nil, nil, err
	}

	return endpoint, data, nil
}

// bffUsage displays the usage of the bff command and its subcommands.
func bffUsage() {
	fmt.Fprintf(os.Stderr, `Backend for Frontend service for chat application
Usage:
    %[1]s [globalflags] bff COMMAND [flags]

COMMAND:
    create-room: Create a new chat room
    history: Get chat room history with enriched user names
    room-list: Get all chat rooms history
    join-room: Creates a new chat room
    invite-room: Creates a new chat room
    stream-chat: Stream chat messages with bidirectional communication
    get-profile: Get current user profile
    update-profile: Update current user profile

Additional help:
    %[1]s bff COMMAND --help
`, os.Args[0])
}
func bffCreateRoomUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] bff create-room -token STRING

Create a new chat room
    -token STRING: 

Example:
    %[1]s bff create-room --token "Aspernatur tenetur libero accusantium laborum."
`, os.Args[0])
}

func bffHistoryUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] bff history -message JSON -token STRING

Get chat room history with enriched user names
    -message JSON: 
    -token STRING: 

Example:
    %[1]s bff history --message '{
      "room_id": "Pariatur non commodi sunt quibusdam."
   }' --token "Totam voluptatum in error suscipit eius impedit."
`, os.Args[0])
}

func bffRoomListUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] bff room-list -token STRING

Get all chat rooms history
    -token STRING: 

Example:
    %[1]s bff room-list --token "Quia illum officiis et."
`, os.Args[0])
}

func bffJoinRoomUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] bff join-room -message JSON -token STRING

Creates a new chat room
    -message JSON: 
    -token STRING: 

Example:
    %[1]s bff join-room --message '{
      "invite_key": "Et quae consequatur expedita."
   }' --token "Sunt doloribus quibusdam nihil sed."
`, os.Args[0])
}

func bffInviteRoomUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] bff invite-room -message JSON -token STRING

Creates a new chat room
    -message JSON: 
    -token STRING: 

Example:
    %[1]s bff invite-room --message '{
      "room_id": "Veritatis error et nulla eius.",
      "user_id": "Sed dolorem."
   }' --token "Mollitia corrupti placeat enim aut."
`, os.Args[0])
}

func bffStreamChatUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] bff stream-chat -token STRING -room-id STRING

Stream chat messages with bidirectional communication
    -token STRING: 
    -room-id STRING: 

Example:
    %[1]s bff stream-chat --token "Autem ipsam officiis rem autem." --room-id "Quo facilis quas voluptatum."
`, os.Args[0])
}

func bffGetProfileUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] bff get-profile -message JSON -token STRING

Get current user profile
    -message JSON: 
    -token STRING: 

Example:
    %[1]s bff get-profile --message '{
      "user_id": "Alias vel."
   }' --token "Nostrum ipsa consequatur vel et inventore."
`, os.Args[0])
}

func bffUpdateProfileUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] bff update-profile -message JSON -token STRING

Update current user profile
    -message JSON: 
    -token STRING: 

Example:
    %[1]s bff update-profile --message '{
      "name": "Rem non sequi rerum dicta autem."
   }' --token "Nostrum quo et quia."
`, os.Args[0])
}
