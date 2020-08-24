package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"io/ioutil"
	"os"
	"time"
	pb "turps/api"
)

const (
	defaultServer = "localhost:50051"
)

var (
	createCmd      = flag.NewFlagSet("changelist", flag.ExitOnError)
	createFromFile = createCmd.String("from_file", "", "input file")
	createServer   = createCmd.String("server", defaultServer, "The address and port of the turpsd server")

	getCmd    = flag.NewFlagSet("get", flag.ExitOnError)
	getServer = getCmd.String("server", defaultServer, "The address and port of the turpsd server")
)

func CreateCommand(args []string) error {
	err := createCmd.Parse(args)
	if err != nil {
		return err
	}
	//fmt.Printf("args %v %d\n", createCmd.Args(), createCmd.NArg())
	if createCmd.NArg() != 1 {
		fmt.Println("Need to specify 1 resource")
		os.Exit(1)
	}
	resourceType := createCmd.Arg(0)
	//fmt.Println("subcommand 'create'")
	//fmt.Println("  resourceType:", resourceType)
	//fmt.Println("  from_file:", *createFromFile)
	//fmt.Println("  server:", *createServer)

	if *createFromFile == "" {
		return errors.New("from_file is required")
	}

	// read file
	dat, err := ioutil.ReadFile(*createFromFile)
	if err != nil {
		return err
	}

	if resourceType == "changelist" {
		err = createChangeList(dat)
	} else if resourceType == "testrun" {
		err = createTestRun(dat)
	}

	return err
}

func createTestRun(dat []byte) error {
	var testRun pb.TestRun
	err := protojson.Unmarshal(dat, &testRun)
	if err != nil {
		return err
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(*createServer, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewTurpsClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.UpsertTestResult(ctx, &pb.UpsertTestRunRequest{TestRun: &testRun})
	if err != nil {
		return err
	}
	fmt.Printf("OK")
	return nil
}

func createChangeList(dat []byte) error {
	var changeList pb.ChangeList
	err := protojson.Unmarshal(dat, &changeList)
	if err != nil {
		return err
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(*createServer, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewTurpsClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.UpsertChangeList(ctx, &pb.UpsertChangeListRequest{ChangeList: &changeList})
	if err != nil {
		return err
	}
	fmt.Printf("OK")
	return nil
}

func GetCommand(args []string) error {
	err := getCmd.Parse(args)
	if err != nil {
		return err
	}
	if getCmd.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Need to specify change_list_id")
		os.Exit(1)
	}
	changeListId := getCmd.Arg(0)

	// Set up a connection to the server.
	conn, err := grpc.Dial(*getServer, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewTurpsClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetChangeList(ctx, &pb.GetChangeListRequest{ChangeListId: changeListId})
	if err != nil {
		return err
	}

	fmt.Println(protojson.Format(r.ChangeList))

	return nil
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("usage: turps <command> [<args>]")
		fmt.Println("The most commonly used commands are: ")
		fmt.Println(" create")
		fmt.Println(" get")
		return
	}

	var err error
	switch os.Args[1] {
	case "create":
		err = CreateCommand(os.Args[2:])
		break
	case "get":
		err = GetCommand(os.Args[2:])
		break
	default:
		err = errors.New("expected create or get subcommand")
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem running subcommand %v", err)
	}

}
