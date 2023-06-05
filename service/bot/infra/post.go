package infra

import (
	"net/url"
	"time"

	"github.com/kapitan123/telegrofler/service/bot/domain"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

type (
	PostModel struct {
		Id                string          `firestore:"-"`
		Type              string          `firestore:"type"`
		OriginalMessageId int             `firestore:"original_message_id,omitempty"`
		ExternalSourceUrl string          `firestore:"external_source_url,omitempty"`
		Reactions         []ReactionModel `firestore:"reactions"`
		PostedOn          time.Time       `firestore:"posted_on"`
		ChatId            int64           `firestore:"chat_id"`
		Poster            UserRefModel    `firestore:"poster"`
	}

	ReactionModel struct {
		ReactToMessageId int          `firestore:"react_to_message_id"`
		Text             string       `firestore:"text"`
		PostedOn         time.Time    `firestore:"posted_on"`
		Reactor          UserRefModel `firestore:"reactor"`
	}

	UserRefModel struct {
		DisplayName string `firestore:"user_name"`
		Id          int64  `firestore:"user_id"`
	}
)

func (pm *PostModel) SetId(id string) {
	pm.Id = id
}

func (urm UserRefModel) toDomainModel() domain.UserRef {
	return domain.UserRef{
		DisplayName: urm.DisplayName,
		Id:          urm.Id,
	}
}

func (rm ReactionModel) toDomainModel() domain.Reaction {
	return domain.Reaction{
		Text:     rm.Text,
		PostedOn: rm.PostedOn,
		Reactor:  rm.Reactor.toDomainModel(),
	}
}

func (pm PostModel) toDomainModel() domain.Post {
	reactions := lo.Map(pm.Reactions, func(rm ReactionModel, _ int) domain.Reaction {
		return rm.toDomainModel()
	})

	url, _ := url.Parse(pm.ExternalSourceUrl)

	mediaType, err := domain.NewMediaTypeFromString(pm.Type)

	if err != nil {
		logrus.Warn(err)
	}

	return domain.Post{
		Id:                pm.Id,
		Type:              mediaType,
		OriginalMessageId: domain.MessageId(pm.OriginalMessageId),
		ExternalSourceUrl: url,
		Reactions:         reactions,
		PostedOn:          pm.PostedOn,
		ChatId:            domain.ChatId(pm.ChatId),
		Poster:            pm.Poster.toDomainModel(),
	}
}

func MapPostToModel(p domain.Post) PostModel {
	reactionModels := lo.Map(p.Reactions, func(r domain.Reaction, _ int) ReactionModel {
		return MapReactionToModel(r)
	})

	return PostModel{
		Id:                p.Id,
		Type:              p.Type.String(),
		ExternalSourceUrl: p.ExternalSourceUrl.String(),
		OriginalMessageId: int(p.OriginalMessageId),
		Reactions:         reactionModels,
		PostedOn:          p.PostedOn,
		ChatId:            int64(p.ChatId),
		Poster:            MapUserRefToModel(p.Poster),
	}
}

func MapUserRefToModel(ur domain.UserRef) UserRefModel {
	return UserRefModel{
		DisplayName: ur.DisplayName,
		Id:          ur.Id,
	}
}

func MapReactionToModel(r domain.Reaction) ReactionModel {
	return ReactionModel{
		Text:     r.Text,
		PostedOn: r.PostedOn,
		Reactor:  MapUserRefToModel(r.Reactor),
	}
}
