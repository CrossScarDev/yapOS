<script setup lang="ts">
import { ref, watch } from 'vue';
import { step, downloadSteps, canContinue } from '../global';
import { DownloadOS, ExtractPlaydateOS } from '../../wailsjs/go/main/App';

const funnyLoader = ref<boolean>();
const yapos = ref(true);
const indexos = ref(false);
const funnyos = ref(false);
const selectedOS = ref<"yapOS" | "Index OS" | "FunnyOS">("yapOS");

const complete = ref(false)
const status = ref('');

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
    await ExtractPlaydateOS();
    if (funnyLoader.value) {
      status.value = "Downloading FunnyLoader...";
      await DownloadOS("FunnyLoader", "https://github.com/RintaDev5792/FunnyLoader/releases/latest/download/FunnyLoader.pdx.zip", "FunnyLoader.*.pdx.zip", "Launcher.pdx");
      status.value = "Extracting FunnyLoader...";
      status.value = "Patching PlaydateOS with FunnyLoader...";
      if (yapos.value) {
        status.value = "Downloading yapOS...";
        await DownloadOS("yapOS", "https://github.com/CrossScarDev/yapOS/releases/latest/download/yapOS.pdx.zip", "yapOS.*.pdx.zip", "yapOS.pdx");
        status.value = "Extracting yapOS...";
        status.value = "Patching PlaydateOS with yapOS...";
      }
      if (indexos.value) {
        status.value = "Downloading Index OS...";
        await DownloadOS("Index OS", "https://github.com/scratchminer/Index-OS/releases/latest/download/IndexOS-Core.pdx.zip", "IndexOS.*.pdx.zip", "IndexOS.pdx");
        status.value = "Extracting Index OS...";
        status.value = "Patching PlaydateOS with Index OS...";
      }
      if (funnyos.value) {
        status.value = "Downloading FunnyOS...";
        await DownloadOS("FunnyOS", "https://github.com/RintaDev5792/FunnyOS/releases/latest/download/FunnyOS.pdx.zip", "FunnyOS.*.pdx.zip", "FunnyOS.pdx");
        status.value = "Extracting FunnyOS...";
        status.value = "Patching PlaydateOS with FunnyOS...";
      }
    } else {
      if (selectedOS.value === "yapOS") {
        status.value = "Downloading yapOS...";
        await DownloadOS("yapOS", "https://github.com/CrossScarDev/yapOS/releases/latest/download/yapOS.pdx.zip", "yapOS.*.pdx.zip", "Launcher.pdx");
        status.value = "Extracting yapOS...";
        status.value = "Patching PlaydateOS with yapOS...";
      } else if (selectedOS.value === "Index OS") {
        status.value = "Downloading Index OS...";
        await DownloadOS("Index OS", "https://github.com/scratchminer/Index-OS/releases/latest/download/IndexOS-Core.pdx.zip", "IndexOS.*.pdx.zip", "Launcher.pdx");
        status.value = "Extracting Index OS...";
        status.value = "Patching PlaydateOS with Index OS...";
      } else if (selectedOS.value === "FunnyOS") {
        status.value = "Downloading FunnyOS...";
        await DownloadOS("FunnyOS", "https://github.com/RintaDev5792/FunnyOS/releases/latest/download/FunnyOS.pdx.zip", "FunnyOS.*.pdx.zip", "Launcher.pdx");
        status.value = "Extracting FunnyOS...";
        status.value = "Patching PlaydateOS with FunnyOS...";
      }
    }

    status.value = "Copying Access Token...";
    status.value = "Repackaging PlaydateOS...";
    status.value = "Uploading Patched PlaydateOS...";
    status.value = "Installing Patched PlaydateOS...";
    status.value = "Cleaning Up...";

    complete.value = true
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
    <template v-if="complete">
      <p>Installation successful.</p>
    </template>
    <template v-else>
      <p>{{ status }}</p>
      <div class="loader" />
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
