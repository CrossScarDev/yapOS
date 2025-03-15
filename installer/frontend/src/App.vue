<script lang="ts" setup>
import { step, downloadSteps, installSteps } from './steps';
import DownloadSteps from './components/DownloadSteps.vue';
</script>
<template>
  <DownloadSteps />
  <div class="controls">
    <button @click="step--">Back</button>
    <button @click="step++">Next</button>
  </div>
</template>
<style>
*,
*::after,
*::before,
*::placeholder {
  margin: 0;
  color: #cdd6f4;
  box-sizing: border-box;
  text-align: center;
  font-family: 'Nunito Variable';
}

#app {
  background: #181825;
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100vh;
  padding: 2rem 2rem 0 2rem;
  gap: 2rem;
  flex-direction: column;
  width: 100vw;
}

button,
input {
  background: #313244;
  border: none;
  padding: 0.75rem 1rem 0.75rem 1rem;
  border-radius: 0.75rem;
  outline: none;

  &:hover {
    background: #45475a;
  }
}

input:focus {
  background: #45475a;
}

.controls {
  display: flex;
  width: 100vw;
  position: absolute;
  bottom: 0;
  padding: 0.75rem;
  justify-content: right;
  gap: 0.75rem;
}

.error {
  color: #f38ba8;
  font-weight: bold;
  font-size: 1.15rem;
}
</style>
