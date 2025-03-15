<script setup lang="ts">
import { step } from '../steps';
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime';
import { ref, watch } from 'vue';

const serialNumber = ref('PDU1-Y');

watch(step, (_newStep, oldStep) => {
  if (oldStep === 2 && !/^PDU1-Y[0-9]{6}$/g.test(serialNumber.value)) step.value--;
});
</script>
<template>
  <template v-if="step === 1">
    <p>
      Please go to https://play.date/devices/, select your Playdate, and remove it from your account. Press next when
      your done.
    </p>
    <button @click="BrowserOpenURL('https://play.date/devices/')">Open in Browser</button>
  </template>
  <template v-else-if="step === 2">
    <p>
      Please input your Playdate's serial number.
    </p>
    <input v-model="serialNumber" placeholder="PDU1-Y" />
  </template>
</template>
