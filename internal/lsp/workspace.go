// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lsp

import (
	"context"

	"golang.org/x/tools/internal/lsp/protocol"
	"golang.org/x/tools/internal/span"
	errors "golang.org/x/xerrors"
)

func (s *Server) changeFolders(ctx context.Context, event protocol.WorkspaceFoldersChangeEvent) error {
	for _, folder := range event.Removed {
		view := s.session.View(folder.Name)
		if view != nil {
			view.Shutdown(ctx)
		} else {
			return errors.Errorf("view %s for %v not found", folder.Name, folder.URI)
		}
	}

	for _, folder := range event.Added {
		if err := s.addView(ctx, folder.Name, span.NewURI(folder.URI)); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) addView(ctx context.Context, name string, uri span.URI) error {
	view := s.session.NewView(ctx, name, uri)
	s.stateMu.Lock()
	state := s.state
	s.stateMu.Unlock()
	options := s.session.Options()
	defer func() { s.session.SetOptions(options) }()
	if state >= serverInitialized {
		s.fetchConfig(ctx, view, &options)
	}
	return nil
}
