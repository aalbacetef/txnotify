<script setup lang="ts">
import { ref } from 'vue';

enum Mode {
  Light = "light",
  Dark = "dark",
};

const currentMode = ref<Mode>(getCurrentMode());

function getCurrentMode(): Mode {
  const defaultMode = Mode.Light;

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

  const v = elem.getAttribute("data-theme", Mode.Light);
  switch (v) {
    case Mode.Light, Mode.Dark:
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
}
</script>

<template>
  <div class="toggle-mode" @click="toggleMode">
    <span>{{ currentMode }}</span>
  </div>
</template>

<style scoped>
.toggle-mode {
  cursor: pointer;
}
</style>
