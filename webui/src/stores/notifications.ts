
import { ref } from 'vue'
import { defineStore } from 'pinia'
import type { Notification } from '@/lib/notifications';
import { Status } from '@/lib/status';


export const useNotificationsStore = defineStore('notifications', () => {
  const notifications = ref<Notification[]>([]);
  const defaultTimeout = 2000;

  function pushNotification(text: string, status: Status = Status.Success) {
    notifications.value.push({
      message: text,
      status,
      at: new Date(),
      duration: defaultTimeout,
    });

    setTimeout(() => notifications.value.shift(), defaultTimeout)
  }

  return { notifications, pushNotification };
})
