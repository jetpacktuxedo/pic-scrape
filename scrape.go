package main

import (
	"context"
	"fmt"
	"os"

	comatproto "github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/atproto/identity"
	"github.com/bluesky-social/indigo/atproto/syntax"
	"github.com/bluesky-social/indigo/xrpc"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please pass a handle as an argument")
		os.Exit(1)
	}
	handle := os.Args[1]
	fmt.Printf("Fetching all media records for %s, this may take a while...\n", handle)
	err := blobDownloadAll(handle)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
}

func blobDownloadAll(handle string) error {
	ctx := context.Background()
	atid, err := syntax.ParseAtIdentifier(handle)
	if err != nil {
		return err
	}

	// resolve the DID document and PDS for this account
	dir := identity.DefaultDirectory()
	ident, err := dir.Lookup(ctx, *atid)
	if err != nil {
		return err
	}

	// create a new API client to connect to the account's PDS
	xrpcc := xrpc.Client{
		Host: ident.PDSEndpoint(),
	}
	if xrpcc.Host == "" {
		return fmt.Errorf("no PDS endpoint for identity")
	}

	topDir := handle
	os.MkdirAll(topDir, os.ModePerm)

	// blob-specific part starts here!
	cursor := ""
	for {
		// loop over batches of CIDs
		resp, err := comatproto.SyncListBlobs(ctx, &xrpcc, cursor, ident.DID.String(), 500, "")
		if err != nil {
			return err
		}
		for _, cidStr := range resp.Cids {
			// if the file already exists, skip
			blobPath := topDir + "/" + cidStr
			if _, err := os.Stat(blobPath); err == nil {
				continue
			}

			// download the entire blob in to memory, then write to disk
			blobBytes, err := comatproto.SyncGetBlob(ctx, &xrpcc, cidStr, ident.DID.String())
			if err != nil {
				fmt.Printf("Warning, failed to download blob: %s", err)
				continue
			}
			if err := os.WriteFile(blobPath, blobBytes, 0666); err != nil {
				fmt.Printf("Warning, failed to write file: %s", err)
				continue
			}
		}

		// a cursor in the result means there are more CIDs to enumerate
		if resp.Cursor != nil && *resp.Cursor != "" {
			cursor = *resp.Cursor
		} else {
			break
		}
	}
	return nil
}
