<template>
  <q-layout view="hHh lpR fFf">

    <q-header reveal elevated class="hidden bg-primary text-white" height-hint="98">
      <q-toolbar>
        <q-toolbar-title>
          <q-avatar>
            <img src="/icons/push.svg">
          </q-avatar>
          {{ title }}
        </q-toolbar-title>
      </q-toolbar>

      <q-tabs>
        <q-route-tab v-for="item in breadcrumbs"
         :key="item.key"
         :to="item.route" :label="item.label" />
        <q-route-tab to="/p/damm/011" label="Yesterday's puzzle" />
        <q-route-tab to="/p/damm/004" label="Today's challenge" />
      </q-tabs>
    </q-header>

    <q-page-container>
      <router-view />
    </q-page-container>

    <q-footer reveal class="hidden bg-grey-8 text-white">
      <q-toolbar>
        <q-toolbar-title>
          <q-avatar>
            <img src="/icons/push.svg">
          </q-avatar>
          <div>Title</div>
        </q-toolbar-title>
      </q-toolbar>
    </q-footer>

  </q-layout>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const title = '[Hexaban Puzzle Title]'

const breadcrumbs = ref([{
  key: 0,
  label: 'ONE',
  route: 'DWS'
}])
</script>
