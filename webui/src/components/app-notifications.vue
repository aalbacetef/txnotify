
<script setup lang="ts">
import {useNotificationsStore} from '@/stores/notifications';
import {Status} from '@/lib/status';

const store = useNotificationsStore();
</script>

<template>
  <div class="notifications-wrapper">
    <div
      v-for="notif in store.notifications"
      :key="notif.message"
      class="notification"
      :class="{
        'is-success': notif.status === Status.Success,
        'is-warn': notif.status === Status.Warn,
        'is-error': notif.status === Status.Error,
      }"
    >
      {{notif.message}}
    </div>
  </div>
</template>

<style  scoped >
.notifications-wrapper {
  display: flex;
  flex-direction: column;
  position: fixed;
  top: 50px;
  right: 50px;
}

.notifications-wrapper .notification {
  color: #fdf9f9;
  font-weight: 500;

  &.is-success{
    background-color: #539953;
  }
}

</style>
