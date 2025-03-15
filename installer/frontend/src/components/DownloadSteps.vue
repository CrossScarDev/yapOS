<script setup lang="ts">
import { step } from '../steps';
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime';
import { ref, watch } from 'vue';
import { GetPin } from '../../wailsjs/go/main/App'

const error = ref('');

const serialNumber = ref('PDU1-Y');
const pin = ref('');

watch(step, async (_newStep, oldStep) => {
  error.value = '';

  if (oldStep === 2) {
    if (!/^PDU1-Y\d{6}$/g.test(serialNumber.value)) {
      step.value--;
      return;
    }
    const pinInfo = await GetPin(serialNumber.value);
    if (pinInfo.pin === undefined) {
      step.value--;
      if (pinInfo.detail !== undefined) error.value = pinInfo.detail;
      return;
    }
    pin.value = pinInfo.pin;
  }
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
    <span class="error" v-if="error != ''">Error: {{ error }}</span>
    <p>
      Please input your Playdate's serial number.
    </p>
    <input v-model="serialNumber" placeholder="PDU1-Y" />
  </template>
  <template v-else-if="step === 3">
    <p>Please go to https://play.date/pin and enter the following pin: <b>{{ pin }}</b></p>
    <button @click="BrowserOpenURL('https://play.date/pin/')">Open in Browser</button>
  </template>
</template>
