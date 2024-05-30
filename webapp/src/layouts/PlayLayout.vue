<template>
  <q-layout view="hHh lpR fFf">

    <q-header reveal elevated class="bg-primary text-white" height-hint="98">
      <q-toolbar>
        <q-toolbar-title>
          <q-avatar>
            <img src="/icons/push.svg">
          </q-avatar>
          Hexaban Puzzle Title
        </q-toolbar-title>
      </q-toolbar>

      <q-tabs>
        <q-route-tab v-for="item in breadcrumbs"
         :key="item.key"
         :to="item.route" :label="item.label" />
        <q-route-tab to="/page2" label="Page Two" />
        <q-route-tab to="/page3" label="Page Three" />
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

const breadcrumbs = ref([{
  key: 0,
  label: 'ONE',
  route: 'DWS'
}])
</script>
