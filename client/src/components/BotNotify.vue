<script setup lang="ts">
import { ref, watch } from 'vue'
import { Switch } from '@headlessui/vue'

const emit = defineEmits<{
  (e: 'updeteBotNotify', value: boolean): boolean
}>()

const newBotNotify = ref(true)

watch(newBotNotify, () => {
  emit('updeteBotNotify', newBotNotify.value)
})
</script>

<template>
  <div class="bot">
    <label class="flex flex-col items-center">
      <p>botの発言を通知する</p>
      <Switch
        v-model="newBotNotify"
        :class="newBotNotify ? 'bg-emerald-400' : 'bg-gray-200'"
        class="relative inline-flex h-[38px] w-[74px] shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus-visible:ring-2 focus-visible:ring-white focus-visible:ring-opacity-75 p-0"
      >
        <span class="sr-only">Use setting</span>
        <span
          aria-hidden="true"
          :class="newBotNotify ? 'translate-x-9' : 'translate-x-0'"
          class="pointer-events-none inline-block h-[34px] w-[34px] transform rounded-full bg-white shadow-lg ring-0 transition duration-200 ease-in-out"
        ></span>
      </Switch>
    </label>
  </div>
</template>
