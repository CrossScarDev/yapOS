<script lang="ts" setup>
import { step, downloadSteps, installSteps, canContinue } from './global';
import DownloadSteps from './components/DownloadSteps.vue';
import InstallSteps from './components/InstallSteps.vue';
</script>
<template>
  <DownloadSteps v-if="step <= downloadSteps" />
  <InstallSteps v-else />
  <div class="controls" v-if="canContinue">
    <button @click="step--" v-if="step != 1 && step != downloadSteps + 1">Back</button>
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
input[type="text"] {
  background: #313244;
  border: none;
  padding: 0.75rem 1rem 0.75rem 1rem;
  border-radius: 0.75rem;
  outline: none;

  &:hover {
    background: #45475a;
  }
}

input[type="text"]:focus {
  background: #45475a;
}

input[type="checkbox"],
input[type="radio"] {
  appearance: none;
  font: inherit;
  color: currentColor;
  border: 0.15rem solid currentColor;
  height: 1rem;
  aspect-ratio: 1;
  display: grid;
  justify-content: center;
  align-items: center;

  &::before {
    content: "";
    width: 0.55rem;
    aspect-ratio: 1;
    scale: 0;
    box-shadow: inset 1rem 1rem #cdd6f4;
  }

  &:checked::before {
    scale: 1;
  }
}

input[type="checkbox"] {
  border-radius: 0.25rem;

  &::before {
    border-radius: 0.05rem;
    transition: 85ms scale ease-in-out;
  }
}

input[type="radio"] {
  border-radius: 100%;

  &::before {
    border-radius: 100%;
    transition: 125ms scale ease-in-out;
  }
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

.loader {
  width: 2.5rem;
  padding: 0.25rem;
  aspect-ratio: 1;
  border-radius: 50%;
  background: #cdd6f4;
  mask: conic-gradient(#0000 10%, #000), linear-gradient(#000 0 0) content-box;
  mask-composite: subtract;
  animation: loading 1s infinite linear;
}

@keyframes loading {
  to {
    transform: rotate(1turn)
  }
}

ul {
  list-style-type: none;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

li {
  display: flex;
  gap: 1rem;
  align-items: center;
}
</style>
