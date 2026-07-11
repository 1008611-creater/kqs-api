import { apiClient } from './client'

export type GroupBuyRoomStatus = 'active' | 'completed' | 'expired' | 'canceled'
export type GroupBuyTargetCount = 3 | 5 | 10

export interface GroupBuyRoom {
  id: number
  plan_id: number
  plan_name: string
  group_id: number
  group_name: string
  owner_user_id: number
  target_count: number
  member_count: number
  bonus_multiplier: number
  base_daily_limit_usd: number
  upgraded_daily_limit_usd: number
  base_weekly_limit_usd: number
  upgraded_weekly_limit_usd: number
  status: GroupBuyRoomStatus
  expires_at: string
  completed_at?: string | null
  created_at: string
  joined: boolean
  can_join: boolean
}

export interface GroupBuyHall {
  rooms: GroupBuyRoom[]
  my_rooms: GroupBuyRoom[]
}

export interface CreateGroupBuyRoomRequest {
  plan_id: number
  target_count: GroupBuyTargetCount
}

export async function listGroupBuys(): Promise<GroupBuyHall> {
  const { data } = await apiClient.get<GroupBuyHall>('/group-buys')
  return data
}

export async function createGroupBuyRoom(payload: CreateGroupBuyRoomRequest): Promise<GroupBuyRoom> {
  const { data } = await apiClient.post<GroupBuyRoom>('/group-buys', payload)
  return data
}

export async function joinGroupBuyRoom(roomID: number): Promise<GroupBuyRoom> {
  const { data } = await apiClient.post<GroupBuyRoom>(`/group-buys/${roomID}/join`)
  return data
}

export const groupBuyAPI = {
  listGroupBuys,
  createGroupBuyRoom,
  joinGroupBuyRoom
}

export default groupBuyAPI
