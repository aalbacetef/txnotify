import type { Status } from "./status";

export type Notification = {
  message: string;
  at: Date;
  duration: number;
  status: Status
}
