package services

import (
	"context"

	"github.com/shellhub-io/shellhub/pkg/models"
	"github.com/shellhub-io/shellhub/pkg/validator"
)

type DeviceTags interface {
	CreateTag(ctx context.Context, uid models.UID, name string) error
	RemoveTag(ctx context.Context, uid models.UID, name string) error
	RenameTag(ctx context.Context, tenantID string, currentName string, newName string) error
	UpdateTag(ctx context.Context, uid models.UID, tags []string) error
	GetTags(ctx context.Context, tenant string) ([]string, int, error)
	DeleteTags(ctx context.Context, tenant string, name string) error
}

// DeviceMaxTags is the number of tags that a device can have.
const DeviceMaxTags = 3

func (s *service) CreateTag(ctx context.Context, uid models.UID, name string) error {
	if !validator.ValidateFieldTag(name) {
		return ErrInvalidFormat
	}

	device, err := s.store.DeviceGet(ctx, uid)
	if err != nil || device == nil {
		return ErrDeviceNotFound
	}

	if len(device.Tags) == DeviceMaxTags {
		return ErrMaxTagReached
	}

	if contains(device.Tags, name) {
		return ErrDuplicateTagName
	}

	return s.store.DeviceCreateTag(ctx, uid, name)
}

func (s *service) RemoveTag(ctx context.Context, uid models.UID, name string) error {
	device, err := s.store.DeviceGet(ctx, uid)
	if err != nil || device == nil {
		return ErrDeviceNotFound
	}

	if !contains(device.Tags, name) {
		return ErrTagNameNotFound
	}

	return s.store.DeviceRemoveTag(ctx, uid, name)
}

func (s *service) RenameTag(ctx context.Context, tenantID string, currentName string, newName string) error {
	if !validator.ValidateFieldTag(newName) {
		return ErrInvalidFormat
	}

	tags, count, err := s.store.DeviceGetTags(ctx, tenantID)
	if err != nil || count == 0 {
		return ErrNoTags
	}

	if !contains(tags, currentName) {
		return ErrTagNameNotFound
	}

	if contains(tags, newName) {
		return ErrDuplicateTagName
	}

	return s.store.DeviceRenameTag(ctx, tenantID, currentName, newName)
}

func (s *service) UpdateTag(ctx context.Context, uid models.UID, tags []string) error {
	listToSet := func(list []string) []string {
		s := make(map[string]bool)
		l := make([]string, 0)
		for _, o := range list {
			if _, ok := s[o]; !ok {
				s[o] = true
				l = append(l, o)
			}
		}

		return l
	}

	if len(tags) > DeviceMaxTags {
		return ErrMaxTagReached
	}

	tagSet := listToSet(tags)

	for _, tag := range tagSet {
		if !validator.ValidateFieldTag(tag) {
			return ErrInvalidFormat
		}
	}

	device, err := s.store.DeviceGet(ctx, uid)
	if err != nil || device == nil {
		return ErrDeviceNotFound
	}

	return s.store.DeviceUpdateTag(ctx, uid, tagSet)
}

func (s *service) GetTags(ctx context.Context, tenant string) ([]string, int, error) {
	namespace, err := s.store.NamespaceGet(ctx, tenant)
	if err != nil || namespace == nil {
		return nil, 0, ErrNamespaceNotFound
	}

	return s.store.DeviceGetTags(ctx, namespace.TenantID)
}

func (s *service) DeleteTags(ctx context.Context, tenant string, name string) error {
	namespace, err := s.store.NamespaceGet(ctx, tenant)
	if err != nil || namespace == nil {
		return ErrNamespaceNotFound
	}

	return s.store.DeviceDeleteTags(ctx, namespace.TenantID, name)
}
