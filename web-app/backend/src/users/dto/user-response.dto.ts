import { User } from "@prisma/client";

export interface UserResponse {
  email: string;
  name: string;
  createdAt: Date;
}

export function mapToUserResponse(user: User): UserResponse {
  return {
    email: user.email,
    name: user.name ?? '',
    createdAt: user.createdAt,
  };
}