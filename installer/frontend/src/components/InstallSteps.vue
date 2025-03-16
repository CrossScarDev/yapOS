<script setup lang="ts">
import { ref, watch } from 'vue';
import { step, downloadSteps, installSteps, canContinue } from '../global';

const funnyLoader = ref(true);
const yapos = ref(true);
const indexos = ref(false);
const funnyos = ref(false);

watch(step, (newStep, oldStep) => {
  if (newStep - downloadSteps === 1) canContinue.value = false;
  if (oldStep - downloadSteps === 1) canContinue.value = true;
});
</script>
<template>
  <template v-if="step - downloadSteps === 1">
    <p>Would you like to use FunnyLoader?</p>
    <div>
      <button @click="funnyLoader = true; step++">Yes</button>
      <button @click="funnyLoader = false; step++">No</button>
    </div>
  </template>
  <template v-else-if="step - downloadSteps === 2">
    <template v-if="funnyLoader">
      <p>Select the operating systems you would like to install:</p>
      <ul>
        <li>
          <input type="checkbox" id="yapos" v-model="yapos" />
          <label for="yapos">yapOS</label>
        </li>
        <li>
          <input type="checkbox" id="indexos" v-model="indexos" />
          <label for="indexos">Index OS</label>
        </li>
        <li>
          <input type="checkbox" id="funnyos" v-model="funnyos" />
          <label for="funnyos">FunnyOS</label>
        </li>
      </ul>
    </template>
    <template v-else>
      <p>Select the operating system you would like to install:</p>
      <ul>
        <li>
          <input type="radio" id="yapos" v-model="yapos" name="os" />
          <label for="yapos">yapOS</label>
        </li>
        <li>
          <input type="radio" id="indexos" v-model="indexos" name="os" />
          <label for="indexos">Index OS</label>
        </li>
        <li>
          <input type="radio" id="funnyos" v-model="funnyos" name="os" />
          <label for="funnyos">FunnyOS</label>
        </li>
      </ul>
    </template>

  </template>
</template>

<style scoped>
div {
  display: flex;
  flex-direction: row;
  gap: 0.75rem;
}
</style>
