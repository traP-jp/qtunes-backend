package model

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	"github.com/hackathon-21-spring-02/back-end/domain"
	"github.com/sapphi-red/go-traq"
	"net/http"
	"strings"
)

func GetComposers(ctx context.Context, accessToken string) ([]*domain.Composer, error) {
	client, auth := newClient(accessToken)
	files, res, err := client.FileApi.GetFiles(auth, &traq.FileApiGetFilesOpts{
		ChannelId: optional.NewInterface(SoundChannelId),
		Limit:     optional.NewInt32(200),
	})
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
	}

	users,res,err:=client.UserApi.GetUsers(auth,nil)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
	}

	postCountByUser:=make(map[string]int)
	fmt.Println(files[100].UploaderId)
	for _, v := range files{
		if strings.HasPrefix(v.Mime, "audio") {
			if _, ok := postCountByUser[*v.UploaderId]; ok {
				postCountByUser[*v.UploaderId]++
			}else {
				postCountByUser[*v.UploaderId]=1
			}
		}
	}
	//idToName:=make(map[*string]string)
	composers:=make([]*domain.Composer,0,len(users))
	for _,user:=range users{
		if val,ok:=postCountByUser[user.Id];ok&&val>0 {
			composers=append(composers,&domain.Composer{
				ID:        user.Id,
				Name:      user.Name,
				PostCount: val,
			})

		}
	}
	fmt.Println(users[0].Name)
	return composers,err
}
