// Copyright 2018 The Nakama Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/heroiclabs/nakama/api"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ApiServer) CreateGroup(ctx context.Context, in *api.CreateGroupRequest) (*api.Group, error) {
	userID := ctx.Value(ctxUserIDKey{}).(uuid.UUID)

	// Before hook.
	if fn := s.runtime.BeforeCreateGroup(); fn != nil {
		beforeFn := func(clientIP, clientPort string) error {
			result, err, code := fn(ctx, s.logger, userID.String(), ctx.Value(ctxUsernameKey{}).(string), ctx.Value(ctxExpiryKey{}).(int64), clientIP, clientPort, in)
			if err != nil {
				return status.Error(code, err.Error())
			}
			if result == nil {
				// If result is nil, requested resource is disabled.
				s.logger.Warn("Intercepted a disabled resource.", zap.Any("resource", ctx.Value(ctxFullMethodKey{}).(string)), zap.String("uid", userID.String()))
				return status.Error(codes.NotFound, "Requested resource was not found.")
			}
			in = result
			return nil
		}

		// Execute the before function lambda wrapped in a trace for stats measurement.
		err := traceApiBefore(ctx, s.logger, ctx.Value(ctxFullMethodKey{}).(string), beforeFn)
		if err != nil {
			return nil, err
		}
	}

	if in.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "Group name must be set.")
	}

	group, err := CreateGroup(ctx, s.logger, s.db, userID, userID, in.GetName(), in.GetLangTag(), in.GetDescription(), in.GetAvatarUrl(), "", in.GetOpen(), -1)
	if err != nil {
		if err == ErrGroupNameInUse {
			return nil, status.Error(codes.InvalidArgument, "Group name is in use.")
		}
		return nil, status.Error(codes.Internal, "Error while trying to create group.")
	}

	// After hook.
	if fn := s.runtime.AfterCreateGroup(); fn != nil {
		afterFn := func(clientIP, clientPort string) {
			fn(ctx, s.logger, userID.String(), ctx.Value(ctxUsernameKey{}).(string), ctx.Value(ctxExpiryKey{}).(int64), clientIP, clientPort, group, in)
		}

		// Execute the after function lambda wrapped in a trace for stats measurement.
		traceApiAfter(ctx, s.logger, ctx.Value(ctxFullMethodKey{}).(string), afterFn)
	}

	return group, nil
}

func (s *ApiServer) UpdateGroup(ctx context.Context, in *api.UpdateGroupRequest) (*empty.Empty, error) {
	userID := ctx.Value(ctxUserIDKey{}).(uuid.UUID)

	// Before hook.
	if fn := s.runtime.BeforeUpdateGroup(); fn != nil {
		beforeFn := func(clientIP, clientPort string) error {
			result, err, code := fn(ctx, s.logger, userID.String(), ctx.Value(ctxUsernameKey{}).(string), ctx.Value(ctxExpiryKey{}).(int64), clientIP, clientPort, in)
			if err != nil {
				return status.Error(code, err.Error())
			}
			if result == nil {
				// If result is nil, requested resource is disabled.
				s.logger.Warn("Intercepted a disabled resource.", zap.Any("resource", ctx.Value(ctxFullMethodKey{}).(string)), zap.String("uid", userID.String()))
				return status.Error(codes.NotFound, "Requested resource was not found.")
			}
			in = result
			return nil
		}

		// Execute the before function lambda wrapped in a trace for stats measurement.
		err := traceApiBefore(ctx, s.logger, ctx.Value(ctxFullMethodKey{}).(string), beforeFn)
		if err != nil {
			return nil, err
		}
	}

	if in.GetGroupId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Group ID must be set.")
	}

	groupID, err := uuid.FromString(in.GetGroupId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Group ID must be a valid ID.")
	}

	if in.GetName() != nil {
		if len(in.GetName().String()) < 1 {
			return nil, status.Error(codes.InvalidArgument, "Group name cannot be empty.")
		}
	}

	if in.GetLangTag() != nil {
		if len(in.GetLangTag().String()) < 1 {
			return nil, status.Error(codes.InvalidArgument, "Group language cannot be empty.")
		}
	}

	err = UpdateGroup(ctx, s.logger, s.db, groupID, userID, uuid.Nil, in.GetName(), in.GetLangTag(), in.GetDescription(), in.GetAvatarUrl(), nil, in.GetOpen(), -1)
	if err != nil {
		if err == ErrGroupPermissionDenied {
			return nil, status.Error(codes.NotFound, "Group not found or you're not allowed to update.")
		} else if err == ErrGroupNoUpdateOps {
			return nil, status.Error(codes.InvalidArgument, "Specify at least one field to update.")
		} else if err == ErrGroupNotUpdated {
			return nil, status.Error(codes.InvalidArgument, "No new fields in group update.")
		}
		return nil, status.Error(codes.Internal, "Error while trying to update group.")
	}

	// After hook.
	if fn := s.runtime.AfterUpdateGroup(); fn != nil {
		afterFn := func(clientIP, clientPort string) {
			fn(ctx, s.logger, userID.String(), ctx.Value(ctxUsernameKey{}).(string), ctx.Value(ctxExpiryKey{}).(int64), clientIP, clientPort, in)
		}

		// Execute the after function lambda wrapped in a trace for stats measurement.
		traceApiAfter(ctx, s.logger, ctx.Value(ctxFullMethodKey{}).(string), afterFn)
	}

	return &empty.Empty{}, nil
}

func (s *ApiServer) DeleteGroup(ctx context.Context, in *api.DeleteGroupRequest) (*empty.Empty, error) {
	userID := ctx.Value(ctxUserIDKey{}).(uuid.UUID)

	// Before hook.
	if fn := s.runtime.BeforeDeleteGroup(); fn != nil {
		beforeFn := func(clientIP, clientPort string) error {
			result, err, code := fn(ctx, s.logger, userID.String(), ctx.Value(ctxUsernameKey{}).(string), ctx.Value(ctxExpiryKey{}).(int64), clientIP, clientPort, in)
			if err != nil {
				return status.Error(code, err.Error())
			}
			if result == nil {
				// If result is nil, requested resource is disabled.
				s.logger.Warn("Intercepted a disabled resource.", zap.Any("resource", ctx.Value(ctxFullMethodKey{}).(string)), zap.String("uid", userID.String()))
				return status.Error(codes.NotFound, "Requested resource was not found.")
			}
			in = result
			return nil
		}

		// Execute the before function lambda wrapped in a trace for stats measurement.
		err := traceApiBefore(ctx, s.logger, ctx.Value(ctxFullMethodKey{}).(string), beforeFn)
		if err != nil {
			return nil, err
		}
	}

	if in.GetGroupId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Group ID must be set.")
	}

	groupID, err := uuid.FromString(in.GetGroupId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Group ID must be a valid ID.")
	}

	err = DeleteGroup(ctx, s.logger, s.db, groupID, userID)
	if err != nil {
		if err == ErrGroupPermissionDenied {
			return nil, status.Error(codes.InvalidArgument, "Group not found or you're not allowed to delete.")
		}
		return nil, status.Error(codes.Internal, "Error while trying to delete group.")
	}

	// After hook.
	if fn := s.runtime.AfterDeleteGroup(); fn != nil {
		afterFn := func(clientIP, clientPort string) {
			fn(ctx, s.logger, userID.String(), ctx.Value(ctxUsernameKey{}).(string), ctx.Value(ctxExpiryKey{}).(int64), clientIP, clientPort, in)
		}

		// Execute the after function lambda wrapped in a trace for stats measurement.
		traceApiAfter(ctx, s.logger, ctx.Value(ctxFullMethodKey{}).(string), afterFn)
	}

	return &empty.Empty{}, nil
}

func (s *ApiServer) JoinGroup(ctx context.Context, in *api.JoinGroupRequest) (*empty.Empty, error) {
	userID := ctx.Value(ctxUserIDKey{}).(uuid.UUID)
	username := ctx.Value(ctxUsernameKey{}).(string)

	// Before hook.
	if fn := s.runtime.BeforeJoinGroup(); fn != nil {
		beforeFn := func(clientIP, clientPort string) error {
			result, err, code := fn(ctx, s.logger, userID.String(), username, ctx.Value(ctxExpiryKey{}).(int64), clientIP, clientPort, in)
			if err != nil {
				return status.Error(code, err.Error())
			}
			if result == nil {
				// If result is nil, requested resource is disabled.
				s.logger.Warn("Intercepted a disabled resource.", zap.Any("resource", ctx.Value(ctxFullMethodKey{}).(string)), zap.String("uid", userID.String()))
				return status.Error(codes.NotFound, "Requested resource was not found.")
			}
			in = result
			return nil
		}

		// Execute the before function lambda wrapped in a trace for stats measurement.
		err := traceApiBefore(ctx, s.logger, ctx.Value(ctxFullMethodKey{}).(string), beforeFn)
		if err != nil {
			return nil, err
		}
	}

	if in.GetGroupId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Group ID must be set.")
	}

	groupID, err := uuid.FromString(in.GetGroupId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Group ID must be a valid ID.")
	}

	err = JoinGroup(ctx, s.logger, s.db, s.router, groupID, userID, username)
	if err != nil {
		if err == ErrGroupNotFound {
			return nil, status.Error(codes.NotFound, "Group not found.")
		} else if err == ErrGroupFull {
			return nil, status.Error(codes.InvalidArgument, "Group is full.")
		}
		return nil, status.Error(codes.Internal, "Error while trying to join group.")
	}

	// After hook.
	if fn := s.runtime.AfterJoinGroup(); fn != nil {
		afterFn := func(clientIP, clientPort string) {
			fn(ctx, s.logger, userID.String(), username, ctx.Value(ctxExpiryKey{}).(int64), clientIP, clientPort, in)
		}

		// Execute the after function lambda wrapped in a trace for stats measurement.
		traceApiAfter(ctx, s.logger, ctx.Value(ctxFullMethodKey{}).(string), afterFn)
	}

	return &empty.Empty{}, nil
}

func (s *ApiServer) LeaveGroup(ctx context.Context, in *api.LeaveGroupRequest) (*empty.Empty, error) {
	userID := ctx.Value(ctxUserIDKey{}).(uuid.UUID)

	// Before hook.
	if fn := s.runtime.BeforeLeaveGroup(); fn != nil {
		beforeFn := func(clientIP, clientPort string) error {
			result, err, code := fn(ctx, s.logger, userID.String(), ctx.Value(ctxUsernameKey{}).(string), ctx.Value(ctxExpiryKey{}).(int64), clientIP, clientPort, in)
			if err != nil {
				return status.Error(code, err.Error())
			}
			if result == nil {
				// If result is nil, requested resource is disabled.
				s.logger.Warn("Intercepted a disabled resource.", zap.Any("resource", ctx.Value(ctxFullMethodKey{}).(string)), zap.String("uid", userID.String()))
				return status.Error(codes.NotFound, "Requested resource was not found.")
			}
			in = result
			return nil
		}

		// Execute the before function lambda wrapped in a trace for stats measurement.
		err := traceApiBefore(ctx, s.logger, ctx.Value(ctxFullMethodKey{}).(string), beforeFn)
		if err != nil {
			return nil, err
		}
	}

	if in.GetGroupId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Group ID must be set.")
	}

	groupID, err := uuid.FromString(in.GetGroupId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Group ID must be a valid ID.")
	}

	err = LeaveGroup(ctx, s.logger, s.db, groupID, userID)
	if err != nil {
		if err == ErrGroupLastSuperadmin {
			return nil, status.Error(codes.InvalidArgument, "Cannot leave group when you are the last superadmin.")
		}
		return nil, status.Error(codes.Internal, "Error while trying to leave group.")
	}

	// After hook.
	if fn := s.runtime.AfterLeaveGroup(); fn != nil {
		afterFn := func(clientIP, clientPort string) {
			fn(ctx, s.logger, userID.String(), ctx.Value(ctxUsernameKey{}).(string), ctx.Value(ctxExpiryKey{}).(int64), clientIP, clientPort, in)
		}

		// Execute the after function lambda wrapped in a trace for stats measurement.
		traceApiAfter(ctx, s.logger, ctx.Value(ctxFullMethodKey{}).(string), afterFn)
	}

	return &empty.Empty{}, nil
}

func (s *ApiServer) AddGroupUsers(ctx context.Context, in *api.AddGroupUsersRequest) (*empty.Empty, error) {
	userID := ctx.Value(ctxUserIDKey{}).(uuid.UUID)

	// Before hook.
	if fn := s.runtime.BeforeAddGroupUsers(); fn != nil {
		beforeFn := func(clientIP, clientPort string) error {
			result, err, code := fn(ctx, s.logger, userID.String(), ctx.Value(ctxUsernameKey{}).(string), ctx.Value(ctxExpiryKey{}).(int64), clientIP, clientPort, in)
			if err != nil {
				return status.Error(code, err.Error())
			}
			if result == nil {
				// If result is nil, requested resource is disabled.
				s.logger.Warn("Intercepted a disabled resource.", zap.Any("resource", ctx.Value(ctxFullMethodKey{}).(string)), zap.String("uid", userID.String()))
				return status.Error(codes.NotFound, "Requested resource was not found.")
			}
			in = result
			return nil
		}

		// Execute the before function lambda wrapped in a trace for stats measurement.
		err := traceApiBefore(ctx, s.logger, ctx.Value(ctxFullMethodKey{}).(string), beforeFn)
		if err != nil {
			return nil, err
		}
	}

	if in.GetGroupId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Group ID must be set.")
	}

	groupID, err := uuid.FromString(in.GetGroupId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Group ID must be a valid ID.")
	}

	if len(in.GetUserIds()) == 0 {
		return &empty.Empty{}, nil
	}

	userIDs := make([]uuid.UUID, 0, len(in.GetUserIds()))
	for _, id := range in.GetUserIds() {
		uid := uuid.FromStringOrNil(id)
		if uid == uuid.Nil {
			return nil, status.Error(codes.InvalidArgument, "User ID must be a valid ID.")
		}
		userIDs = append(userIDs, uid)
	}

	err = AddGroupUsers(ctx, s.logger, s.db, s.router, userID, groupID, userIDs)
	if err != nil {
		if err == ErrGroupPermissionDenied {
			return nil, status.Error(codes.NotFound, "Group not found or permission denied.")
		} else if err == ErrGroupFull {
			return nil, status.Error(codes.InvalidArgument, "Group is full.")
		}
		return nil, status.Error(codes.Internal, "Error while trying to add users to a group.")
	}

	// After hook.
	if fn := s.runtime.AfterAddGroupUsers(); fn != nil {
		afterFn := func(clientIP, clientPort string) {
			fn(ctx, s.logger, userID.String(), ctx.Value(ctxUsernameKey{}).(string), ctx.Value(ctxExpiryKey{}).(int64), clientIP, clientPort, in)
		}

		// Execute the after function lambda wrapped in a trace for stats measurement.
		traceApiAfter(ctx, s.logger, ctx.Value(ctxFullMethodKey{}).(string), afterFn)
	}

	return &empty.Empty{}, nil
}

func (s *ApiServer) KickGroupUsers(ctx context.Context, in *api.KickGroupUsersRequest) (*empty.Empty, error) {
	userID := ctx.Value(ctxUserIDKey{}).(uuid.UUID)

	// Before hook.
	if fn := s.runtime.BeforeKickGroupUsers(); fn != nil {
		beforeFn := func(clientIP, clientPort string) error {
			result, err, code := fn(ctx, s.logger, userID.String(), ctx.Value(ctxUsernameKey{}).(string), ctx.Value(ctxExpiryKey{}).(int64), clientIP, clientPort, in)
			if err != nil {
				return status.Error(code, err.Error())
			}
			if result == nil {
				// If result is nil, requested resource is disabled.
				s.logger.Warn("Intercepted a disabled resource.", zap.Any("resource", ctx.Value(ctxFullMethodKey{}).(string)), zap.String("uid", userID.String()))
				return status.Error(codes.NotFound, "Requested resource was not found.")
			}
			in = result
			return nil
		}

		// Execute the before function lambda wrapped in a trace for stats measurement.
		err := traceApiBefore(ctx, s.logger, ctx.Value(ctxFullMethodKey{}).(string), beforeFn)
		if err != nil {
			return nil, err
		}
	}

	if in.GetGroupId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Group ID must be set.")
	}

	groupID, err := uuid.FromString(in.GetGroupId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Group ID must be a valid ID.")
	}

	if len(in.GetUserIds()) == 0 {
		return &empty.Empty{}, nil
	}

	userIDs := make([]uuid.UUID, 0, len(in.GetUserIds()))
	for _, id := range in.GetUserIds() {
		uid := uuid.FromStringOrNil(id)
		if uid == uuid.Nil {
			return nil, status.Error(codes.InvalidArgument, "User ID must be a valid ID.")
		}
		userIDs = append(userIDs, uid)
	}

	if err = KickGroupUsers(ctx, s.logger, s.db, userID, groupID, userIDs); err != nil {
		if err == ErrGroupPermissionDenied {
			return nil, status.Error(codes.NotFound, "Group not found or permission denied.")
		}
		return nil, status.Error(codes.Internal, "Error while trying to kick users from a group.")
	}

	// After hook.
	if fn := s.runtime.AfterKickGroupUsers(); fn != nil {
		afterFn := func(clientIP, clientPort string) {
			fn(ctx, s.logger, userID.String(), ctx.Value(ctxUsernameKey{}).(string), ctx.Value(ctxExpiryKey{}).(int64), clientIP, clientPort, in)
		}

		// Execute the after function lambda wrapped in a trace for stats measurement.
		traceApiAfter(ctx, s.logger, ctx.Value(ctxFullMethodKey{}).(string), afterFn)
	}

	return &empty.Empty{}, nil
}

func (s *ApiServer) PromoteGroupUsers(ctx context.Context, in *api.PromoteGroupUsersRequest) (*empty.Empty, error) {
	userID := ctx.Value(ctxUserIDKey{}).(uuid.UUID)

	// Before hook.
	if fn := s.runtime.BeforePromoteGroupUsers(); fn != nil {
		beforeFn := func(clientIP, clientPort string) error {
			result, err, code := fn(ctx, s.logger, userID.String(), ctx.Value(ctxUsernameKey{}).(string), ctx.Value(ctxExpiryKey{}).(int64), clientIP, clientPort, in)
			if err != nil {
				return status.Error(code, err.Error())
			}
			if result == nil {
				// If result is nil, requested resource is disabled.
				s.logger.Warn("Intercepted a disabled resource.", zap.Any("resource", ctx.Value(ctxFullMethodKey{}).(string)), zap.String("uid", userID.String()))
				return status.Error(codes.NotFound, "Requested resource was not found.")
			}
			in = result
			return nil
		}

		// Execute the before function lambda wrapped in a trace for stats measurement.
		err := traceApiBefore(ctx, s.logger, ctx.Value(ctxFullMethodKey{}).(string), beforeFn)
		if err != nil {
			return nil, err
		}
	}

	if in.GetGroupId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Group ID must be set.")
	}

	groupID, err := uuid.FromString(in.GetGroupId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Group ID must be a valid ID.")
	}

	if len(in.GetUserIds()) == 0 {
		return &empty.Empty{}, nil
	}

	userIDs := make([]uuid.UUID, 0, len(in.GetUserIds()))
	for _, id := range in.GetUserIds() {
		uid := uuid.FromStringOrNil(id)
		if uid == uuid.Nil {
			return nil, status.Error(codes.InvalidArgument, "User ID must be a valid ID.")
		}
		userIDs = append(userIDs, uid)
	}

	err = PromoteGroupUsers(ctx, s.logger, s.db, userID, groupID, userIDs)
	if err != nil {
		if err == ErrGroupPermissionDenied {
			return nil, status.Error(codes.NotFound, "Group not found or permission denied.")
		} else if err == ErrGroupFull {
			return nil, status.Error(codes.InvalidArgument, "Group is full.")
		}
		return nil, status.Error(codes.Internal, "Error while trying to promote users in a group.")
	}

	// After hook.
	if fn := s.runtime.AfterPromoteGroupUsers(); fn != nil {
		afterFn := func(clientIP, clientPort string) {
			fn(ctx, s.logger, userID.String(), ctx.Value(ctxUsernameKey{}).(string), ctx.Value(ctxExpiryKey{}).(int64), clientIP, clientPort, in)
		}

		// Execute the after function lambda wrapped in a trace for stats measurement.
		traceApiAfter(ctx, s.logger, ctx.Value(ctxFullMethodKey{}).(string), afterFn)
	}

	return &empty.Empty{}, nil
}

func (s *ApiServer) ListGroupUsers(ctx context.Context, in *api.ListGroupUsersRequest) (*api.GroupUserList, error) {
	// Before hook.
	if fn := s.runtime.BeforeListGroupUsers(); fn != nil {
		beforeFn := func(clientIP, clientPort string) error {
			result, err, code := fn(ctx, s.logger, ctx.Value(ctxUserIDKey{}).(uuid.UUID).String(), ctx.Value(ctxUsernameKey{}).(string), ctx.Value(ctxExpiryKey{}).(int64), clientIP, clientPort, in)
			if err != nil {
				return status.Error(code, err.Error())
			}
			if result == nil {
				// If result is nil, requested resource is disabled.
				s.logger.Warn("Intercepted a disabled resource.", zap.Any("resource", ctx.Value(ctxFullMethodKey{}).(string)), zap.String("uid", ctx.Value(ctxUserIDKey{}).(uuid.UUID).String()))
				return status.Error(codes.NotFound, "Requested resource was not found.")
			}
			in = result
			return nil
		}

		// Execute the before function lambda wrapped in a trace for stats measurement.
		err := traceApiBefore(ctx, s.logger, ctx.Value(ctxFullMethodKey{}).(string), beforeFn)
		if err != nil {
			return nil, err
		}
	}

	if in.GetGroupId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Group ID must be set.")
	}

	groupID, err := uuid.FromString(in.GetGroupId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Group ID must be a valid ID.")
	}

	groupUsers, err := ListGroupUsers(ctx, s.logger, s.db, s.tracker, groupID)
	if err != nil {
		return nil, status.Error(codes.Internal, "Error while trying to list users in a group.")
	}

	// After hook.
	if fn := s.runtime.AfterListGroupUsers(); fn != nil {
		afterFn := func(clientIP, clientPort string) {
			fn(ctx, s.logger, ctx.Value(ctxUserIDKey{}).(uuid.UUID).String(), ctx.Value(ctxUsernameKey{}).(string), ctx.Value(ctxExpiryKey{}).(int64), clientIP, clientPort, groupUsers, in)
		}

		// Execute the after function lambda wrapped in a trace for stats measurement.
		traceApiAfter(ctx, s.logger, ctx.Value(ctxFullMethodKey{}).(string), afterFn)
	}

	return groupUsers, nil
}

func (s *ApiServer) ListUserGroups(ctx context.Context, in *api.ListUserGroupsRequest) (*api.UserGroupList, error) {
	// Before hook.
	if fn := s.runtime.BeforeListUserGroups(); fn != nil {
		beforeFn := func(clientIP, clientPort string) error {
			result, err, code := fn(ctx, s.logger, ctx.Value(ctxUserIDKey{}).(uuid.UUID).String(), ctx.Value(ctxUsernameKey{}).(string), ctx.Value(ctxExpiryKey{}).(int64), clientIP, clientPort, in)
			if err != nil {
				return status.Error(code, err.Error())
			}
			if result == nil {
				// If result is nil, requested resource is disabled.
				s.logger.Warn("Intercepted a disabled resource.", zap.Any("resource", ctx.Value(ctxFullMethodKey{}).(string)), zap.String("uid", ctx.Value(ctxUserIDKey{}).(uuid.UUID).String()))
				return status.Error(codes.NotFound, "Requested resource was not found.")
			}
			in = result
			return nil
		}

		// Execute the before function lambda wrapped in a trace for stats measurement.
		err := traceApiBefore(ctx, s.logger, ctx.Value(ctxFullMethodKey{}).(string), beforeFn)
		if err != nil {
			return nil, err
		}
	}

	if in.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "User ID must be set.")
	}

	userID, err := uuid.FromString(in.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Group ID must be a valid ID.")
	}

	userGroups, err := ListUserGroups(ctx, s.logger, s.db, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "Error while trying to list groups for a user.")
	}

	// After hook.
	if fn := s.runtime.AfterListUserGroups(); fn != nil {
		afterFn := func(clientIP, clientPort string) {
			fn(ctx, s.logger, ctx.Value(ctxUserIDKey{}).(uuid.UUID).String(), ctx.Value(ctxUsernameKey{}).(string), ctx.Value(ctxExpiryKey{}).(int64), clientIP, clientPort, userGroups, in)
		}

		// Execute the after function lambda wrapped in a trace for stats measurement.
		traceApiAfter(ctx, s.logger, ctx.Value(ctxFullMethodKey{}).(string), afterFn)
	}

	return userGroups, nil
}

func (s *ApiServer) ListGroups(ctx context.Context, in *api.ListGroupsRequest) (*api.GroupList, error) {
	// Before hook.
	if fn := s.runtime.BeforeListGroups(); fn != nil {
		beforeFn := func(clientIP, clientPort string) error {
			result, err, code := fn(ctx, s.logger, ctx.Value(ctxUserIDKey{}).(uuid.UUID).String(), ctx.Value(ctxUsernameKey{}).(string), ctx.Value(ctxExpiryKey{}).(int64), clientIP, clientPort, in)
			if err != nil {
				return status.Error(code, err.Error())
			}
			if result == nil {
				// If result is nil, requested resource is disabled.
				s.logger.Warn("Intercepted a disabled resource.", zap.Any("resource", ctx.Value(ctxFullMethodKey{}).(string)), zap.String("uid", ctx.Value(ctxUserIDKey{}).(uuid.UUID).String()))
				return status.Error(codes.NotFound, "Requested resource was not found.")
			}
			in = result
			return nil
		}

		// Execute the before function lambda wrapped in a trace for stats measurement.
		err := traceApiBefore(ctx, s.logger, ctx.Value(ctxFullMethodKey{}).(string), beforeFn)
		if err != nil {
			return nil, err
		}
	}

	limit := 1
	if in.GetLimit() != nil {
		if in.GetLimit().Value < 1 || in.GetLimit().Value > 100 {
			return nil, status.Error(codes.InvalidArgument, "Invalid limit - limit must be between 1 and 100.")
		}
		limit = int(in.GetLimit().Value)
	}

	omit_full := false
	if in.GetOmitFull() != nil {
		omit_full = in.GetOmitFull().Value
	}

	groups, err := ListGroups(ctx, s.logger, s.db, in.GetName(), limit, in.GetCursor(), omit_full)
	if err != nil {
		return nil, status.Error(codes.Internal, "Error while trying to list groups.")
	}

	// After hook.
	if fn := s.runtime.AfterListGroups(); fn != nil {
		afterFn := func(clientIP, clientPort string) {
			fn(ctx, s.logger, ctx.Value(ctxUserIDKey{}).(uuid.UUID).String(), ctx.Value(ctxUsernameKey{}).(string), ctx.Value(ctxExpiryKey{}).(int64), clientIP, clientPort, groups, in)
		}

		// Execute the after function lambda wrapped in a trace for stats measurement.
		traceApiAfter(ctx, s.logger, ctx.Value(ctxFullMethodKey{}).(string), afterFn)
	}

	return groups, nil
}
