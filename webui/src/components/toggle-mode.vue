<script setup lang="ts">
import { onMounted } from 'vue';
import { ref } from 'vue';

enum Mode {
  Light = "light",
  Dark = "dark",
};

const currentMode = ref<Mode>(getCurrentMode());

onMounted(() => {
  setMode(currentMode.value);
})

function getCurrentMode(): Mode {
  const defaultMode = Mode.Light;

  const localVal = localStorage.getItem("theme");
  if (localVal !== null) {
    switch (localVal) {
      case Mode.Light:
      case Mode.Dark:
        return localVal as Mode;
    }
  }

  if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
    return Mode.Dark;
  }

  if (window.matchMedia && window.matchMedia('(prefers-color-scheme: light)').matches) {
    return Mode.Light;
  }

  const elem = document.querySelector('html');
  if (elem === null) {
    throw new Error("html element is null");
  }

  const v = elem.getAttribute("data-theme");
  switch (v) {
    case Mode.Light:
    case Mode.Dark:
      return v as Mode;
    default:
      return defaultMode;
  }
}

function toggleMode() {
  if (currentMode.value === Mode.Light) {
    currentMode.value = Mode.Dark;
    setMode(currentMode.value);
    return;
  }

  currentMode.value = Mode.Light;
  setMode(currentMode.value);
}

function setMode(mode: Mode): void {
  const elem = document.querySelector('html');
  if (elem === null) {
    throw new Error("could not get html element");
  }

  elem.setAttribute("data-theme", mode);
  localStorage.setItem("theme", mode);
}
</script>

<template>
  <button class="button toggle-mode" @click="toggleMode">
    <span>{{ currentMode }}</span>
  </button>
</template>

<style scoped>
.toggle-mode {
  cursor: pointer;
  width: 150px;
  padding: 10px 15px !important;
}
</style>
