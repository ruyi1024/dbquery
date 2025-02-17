/*
Copyright 2014-2022 The Lepus Team Group, website: https://www.lepus.cc
Licensed under the GNU General Public License, Version 3.0 (the "GPLv3 License");
You may not use this file except in compliance with the License.
You may obtain a copy of the License at
    https://www.gnu.org/licenses/gpl-3.0.html
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
Special note:
Please do not use this source code for any commercial purpose,
or use it for commercial purposes after secondary development, otherwise you may bear legal risks.
*/

package mongodb

import (
	"context"
	"fmt"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(host, port, username, password, database string) (*mongo.Client, error) {
	var mongoUrl string
	if username != "" && password != "" {
		mongoUrl = fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
	} else {
		mongoUrl = fmt.Sprintf("mongodb://%s:%s", host, port)
	}

	//连接到mongo
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUrl))
	if err != nil {
		return nil, err
	}

	//检测连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	// defer func() {
	// 	if err := client.Disconnect(context.TODO()); err != nil {
	// 		return nil, err
	// }()
	return client, nil
}

func QConnect(host, port, username, password, database string) (*qmgo.Client, error) {
	var mongoUrl string
	if username != "" && password != "" {
		mongoUrl = fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
	} else {
		mongoUrl = fmt.Sprintf("mongodb://%s:%s", host, port)
	}
	ctx := context.Background()
	client, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: mongoUrl})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func ListDatabase(client *mongo.Client) ([]string, error) {
	var databases []string
	// Use a filter to only select non-empty databases.
	databases, err := client.ListDatabaseNames(context.TODO(),
		bson.D{})
	//bson.D{{"empty", false}})
	return databases, err
}

func ListCollection(client *mongo.Client, database string) ([]string, error) {
	var collections []string
	// Use a filter to only select non-empty databases.
	collections, err := client.Database(database).ListCollectionNames(context.TODO(), bson.D{})
	return collections, err
}
