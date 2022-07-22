package service

import (
	"github.com/mashingan/smapping"
	"gitlab.zalopay.vn/top/intern/vybnt/gallery-backend/gallery/dto"
	"gitlab.zalopay.vn/top/intern/vybnt/gallery-backend/gallery/entity"
	"gitlab.zalopay.vn/top/intern/vybnt/gallery-backend/gallery/repository"
	"log"
)

type LikeService interface {
	Like(like dto.LikeDTO) entity.Like
	Unlike(like entity.Like)
	AllLike(postID uint64) []entity.Like
	CountLike(postID uint64) int
}

type likeService struct {
	likeRepository repository.LikeRepository
}

func NewLikeService(likeRepo repository.LikeRepository) LikeService {
	return &likeService{
		likeRepository: likeRepo,
	}
}

func (service *likeService) Like(like dto.LikeDTO) entity.Like {
	likeToSave := entity.Like{}
	err := smapping.FillStruct(&likeToSave, smapping.MapFields(&like))
	if err != nil {
		log.Fatalf("Failed map %v", err)
	}
	savedLike := service.likeRepository.Like(likeToSave)
	return savedLike
}

func (service *likeService) Unlike(like entity.Like) {
	likeToDelete := entity.Like{}
	err := smapping.FillStruct(&likeToDelete, smapping.MapFields(&like))
	if err != nil {
		log.Fatalf("Failed map %v", err)
	}
	service.likeRepository.Unlike(likeToDelete)
}

func (service *likeService) AllLike(postID uint64) []entity.Like {
	return service.likeRepository.AllLikes(postID)
}

func (service *likeService) CountLike(postID uint64) int {
	return service.likeRepository.CountLikes(postID)
}
