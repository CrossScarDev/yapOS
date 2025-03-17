<script setup lang="ts">
import { ref, watch } from 'vue';
import { step, downloadSteps, canContinue } from '../global';
import { CompressPlaydateOS, DownloadOS, ExtractOS, ExtractPlaydateOS, GetSerialPorts, InstallPatchedPlaydateOS, UploadPatchedPlaydateOS, CleanUp } from '../../wailsjs/go/main/App';

const funnyLoader = ref<boolean>();
const yapos = ref(true);
const indexos = ref(false);
const funnyos = ref(false);
const selectedOS = ref<"yapOS" | "Index OS" | "FunnyOS">("yapOS");

const status = ref('');

const ports = ref<string[]>([]);
const selectedPort = ref<string>();

watch(step, async (newStep, oldStep) => {
  if (newStep - downloadSteps === 1) {
    canContinue.value = false;
    return;
  }
  if (oldStep - downloadSteps === 1) {
    canContinue.value = true;
    return;
  }
  if (oldStep - downloadSteps === 2 && newStep - downloadSteps === 3) {
    canContinue.value = false;
    status.value = "Extracting PlaydateOS...";
    await ExtractPlaydateOS(funnyLoader.value as boolean);
    if (funnyLoader.value) {
      status.value = "Downloading FunnyLoader...";
      await DownloadOS("FunnyLoader", "https://github.com/RintaDev5792/FunnyLoader/releases/latest/download/FunnyLoader.pdx.zip", "FunnyLoader.*.pdx.zip", "Launcher.pdx");
      status.value = "Extracting FunnyLoader...";
      await ExtractOS("FunnyLoader", "FunnyLoader.pdx");
      if (yapos.value) {
        status.value = "Downloading yapOS...";
        await DownloadOS("yapOS", "https://github.com/CrossScarDev/yapOS/releases/latest/download/yapOS.pdx.zip", "yapOS.*.pdx.zip", "yapOS.pdx");
        status.value = "Extracting yapOS...";
        await ExtractOS("yapOS", "yapOS.pdx");
      }
      if (indexos.value) {
        status.value = "Downloading Index OS...";
        await DownloadOS("Index OS", "https://github.com/scratchminer/Index-OS/releases/latest/download/IndexOS-Core.pdx.zip", "IndexOS.*.pdx.zip", "IndexOS.pdx");
        status.value = "Extracting Index OS...";
        await ExtractOS("Index OS", "IndexOS-Core.pdx");
      }
      if (funnyos.value) {
        status.value = "Downloading FunnyOS...";
        await DownloadOS("FunnyOS", "https://github.com/RintaDev5792/FunnyOS/releases/latest/download/FunnyOS.pdx.zip", "FunnyOS.*.pdx.zip", "FunnyOS.pdx");
        status.value = "Extracting FunnyOS...";
        await ExtractOS("FunnyOS", "FunnyOS.pdx");
      }
    } else {
      if (selectedOS.value === "yapOS") {
        status.value = "Downloading yapOS...";
        await DownloadOS("yapOS", "https://github.com/CrossScarDev/yapOS/releases/latest/download/yapOS.pdx.zip", "yapOS.*.pdx.zip", "Launcher.pdx");
        status.value = "Extracting yapOS...";
        await ExtractOS("yapOS", "yapOS.pdx");
      } else if (selectedOS.value === "Index OS") {
        status.value = "Downloading Index OS...";
        await DownloadOS("Index OS", "https://github.com/scratchminer/Index-OS/releases/latest/download/IndexOS-Core.pdx.zip", "IndexOS.*.pdx.zip", "Launcher.pdx");
        status.value = "Extracting Index OS...";
        await ExtractOS("Index OS", "IndexOS-Core.pdx");
      } else if (selectedOS.value === "FunnyOS") {
        status.value = "Downloading FunnyOS...";
        await DownloadOS("FunnyOS", "https://github.com/RintaDev5792/FunnyOS/releases/latest/download/FunnyOS.pdx.zip", "FunnyOS.*.pdx.zip", "Launcher.pdx");
        status.value = "Extracting FunnyOS...";
        await ExtractOS("FunnyOS", "FunnyOS.pdx");
      }
    }

    status.value = "Repackaging PlaydateOS...";
    await CompressPlaydateOS()
    status.value = "Fetching List of Serial Ports...";
    ports.value = await GetSerialPorts();

    step.value++;
    return;
  }
  if (oldStep - downloadSteps === 4 && newStep - downloadSteps === 5) {
    status.value = "Uploading Patched PlaydateOS...";
    await UploadPatchedPlaydateOS(selectedPort.value as string);
    step.value++;
    canContinue.value = true;
    return;
  }
  if (oldStep - downloadSteps === 6 && newStep - downloadSteps === 7) {
    await InstallPatchedPlaydateOS(selectedPort.value as string);
    return;
  }
  if (oldStep - downloadSteps === 7 && newStep - downloadSteps === 8) {
    await CleanUp(selectedPort.value as string);
    step.value++;
  }
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
          <input type="radio" id="yapos" value="yapOS" v-model="selectedOS" name="os" />
          <label for="yapos">yapOS</label>
        </li>
        <li>
          <input type="radio" id="indexos" value="Index OS" v-model="selectedOS" name="os" />
          <label for="indexos">Index OS</label>
        </li>
        <li>
          <input type="radio" id="funnyos" value="FunnyOS" v-model="selectedOS" name="os" />
          <label for="funnyos">FunnyOS</label>
        </li>
      </ul>
    </template>
  </template>
  <template v-else-if="step - downloadSteps === 3">
    <p>{{ status }}</p>
    <div class="loader" />
  </template>
  <template v-else-if="step - downloadSteps === 4">
    <p>Please select your Playdate's serial port.</p>
    <button v-for="(port, i) in ports" :key="i" @click="selectedPort = port; step++">{{ port }}</button>
    <a href="#" @click="async () => ports = await GetSerialPorts()">Reload Ports</a>
  </template>
  <template v-else-if="step - downloadSteps === 5">
    <p>{{ status }}</p>
    <div class="loader" />
  </template>
  <template v-else-if="step - downloadSteps === 6">
    <p>
      If your Playdate hasn't returned to the Launcher yet, please press <b>A</b> and wait until it returns to the
      launcher.
    </p>
  </template>
  <template v-else-if="step - downloadSteps === 7">
    <p>
      In a few seconds, your Playdate has begin to install your chosen operating systems, press next when it is
      finished.
    </p>
  </template>
  <template v-else-if="step - downloadSteps === 8">
    <p>Cleaning Up...</p>
    <div class="loader" />
  </template>
</template>

<style scoped>
div {
  display: flex;
  flex-direction: row;
  gap: 0.75rem;
}
</style>
