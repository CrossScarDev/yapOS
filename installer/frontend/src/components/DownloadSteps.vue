<script setup lang="ts">
import { step } from '../steps';
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime';
import { ref, watch } from 'vue';
import { GetPin, FinishRegistration, DownloadOS } from '../../wailsjs/go/main/App'

const error = ref('');

const serialNumber = ref('PDU1-Y');
const pin = ref('');

let accessToken: string;

watch(step, async (newStep, oldStep) => {
  error.value = '';

  if (oldStep === 2 && newStep === 3) {
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
    return;
  }
  if (oldStep === 3 && newStep === 4) {
    const info = await FinishRegistration(serialNumber.value);
    if (info["access_token"] === undefined) {
      step.value--;
      if (info.detail !== undefined) error.value = info.detail;
      if (info.registered !== undefined && !info.registered) error.value = "The pin was not entered."
      return;
    }
    accessToken = info["access_token"];

    await DownloadOS(accessToken);
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
  <template v-else-if="step === 4">
    <p>Downloading PlaydateOS...</p>
    <div class="loader" />
  </template>
</template>
