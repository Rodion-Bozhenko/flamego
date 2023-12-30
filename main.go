// Flamego - GUI for firestore
package main

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	g "github.com/AllenDang/giu"
	"github.com/sqweek/dialog"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type doc = map[string]interface{}

var (
	leftMenuWidth      float32 = 200
	selectedCollection string
	app                *firebase.App
	firebaseCtx        = context.Background()
	collectionsRef     []*firestore.CollectionRef
	client             *firestore.Client
	docsIter           *firestore.DocumentIterator
	columns            []*g.TableColumnWidget
	rows               []*g.TableRowWidget
	docs               []doc
	editingDoc         doc
)

func main() {
	w := g.NewMasterWindow("Flamego", 1000, 700, 0)
	w.Run(loop)
}

func loop() {
	g.SingleWindow().Layout(
		g.Popup("editing-doc-modal").Layout(g.Label("This is a popup"), g.Button("Close").OnClick(func() { g.CloseCurrentPopup() })),
		g.Buttonf("Click").OnClick(func() { g.OpenPopup("editing-doc-modal") }),
		g.SplitLayout(g.DirectionVertical, &leftMenuWidth,
			g.Layout{
				g.Label("Collections"),
				g.Column(renderCollectionButtons(collectionsRef)...),
			},
			g.Layout{
				g.Label("Main Frame"),
				g.Button("Add service account file").OnClick(func() {
					path, err := promptServiceAccountPath()
					if err != nil {
						g.Msgbox("Something went wrong", fmt.Sprintf("Cannot load file: %s", err.Error()))
						return
					}
					app, err = logApp(path)
					if err != nil {
						g.Msgbox("Error", fmt.Sprintf("Cannot initialize app: %s\n", err.Error()))
					}

					client, err = app.Firestore(firebaseCtx)
					if err != nil {
						g.Msgbox("Firestore Error", fmt.Sprintf("Cannot connect to Firestore: %s\n", err.Error()))
					}

					collectionsRef, err = client.Collections(firebaseCtx).GetAll()
					if err != nil {
						g.Msgbox("Firestore Error", fmt.Sprintf("Cannot retrieve all collections: %s\n", err.Error()))
					}
				}),
				renderDocsTable(docsIter),
				g.PrepareMsgbox(),
			},
		),
	)
}

func renderCollectionButtons(colls []*firestore.CollectionRef) []g.Widget {
	buttons := make([]g.Widget, 0, len(colls))
	for _, coll := range colls {
		buttons = append(buttons, newCollectionButton(coll.ID))
	}

	return buttons
}

func renderDocsTable(docsIter *firestore.DocumentIterator) g.Widget {
	if docsIter == nil {
		return g.Label("")
	}

	if len(docs) == 0 {
		for {
			docSnap, err := docsIter.Next()
			if errors.Is(err, iterator.Done) {
				docsIter.Stop()
				break
			}
			if err != nil {
				docsIter.Stop()
				panic(err.Error())
			}
			data := docSnap.Data()
			if data != nil {
				docs = append(docs, data)
			}
		}
	}

	var keys []string

	for key := range docs[0] {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	if len(columns) == 0 {
		for _, key := range keys {
			columns = append(columns, g.TableColumn(key))
		}
	}

	if len(rows) == 0 {
		for _, doc := range docs {
			var rowCells []g.Widget
			for _, k := range keys {
				rowCells = append(rowCells, g.Buttonf("%v", doc[k]).OnClick(func() {
					editingDoc = doc
					g.OpenPopup("editing-doc-modal")
				}))
			}
			rows = append(rows, g.TableRow(rowCells...))
		}
	}

	return g.Table().FastMode(true).Columns(columns...).Rows(rows...)
}

func renderDocPopup(doc doc) g.Widget {
	return g.Popup("editing-doc-modal").Layout(g.Label("This is a popup"), g.Button("Close").OnClick(func() { g.CloseCurrentPopup() }))
}

func getAllDocs(client *firestore.Client, collID string) *firestore.DocumentIterator {
	if client != nil {
		return client.Collection(collID).Documents(firebaseCtx)
	}
	return nil
}

func newCollectionButton(title string) *g.ButtonWidget {
	return g.Button(title).OnClick(func() {
		selectedCollection = title
		docsIter = getAllDocs(client, title)
		docs = []map[string]interface{}{}
		columns = []*g.TableColumnWidget{}
		rows = []*g.TableRowWidget{}
	})
}

func promptServiceAccountPath() (string, error) {
	path, err := dialog.File().Filter("Service account json file", "json").Load()
	if err != nil {
		if err == dialog.ErrCancelled {
			return "", nil
		}
		return "", err
	}
	return path, nil
}

func logApp(path string) (*firebase.App, error) {
	opt := option.WithCredentialsFile(path)
	app, err := firebase.NewApp(firebaseCtx, nil, opt)
	if err != nil {
		return nil, err
	}

	return app, nil
}
