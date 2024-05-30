<!--
  Copyright (c) 2024 Symbol Not Found L.L.C.

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

  github:SymbolNotFound/hexoban/webapp/src/components/HexGrid.vue
-->

<template>
  <!-- TODO should be computing viewBox based on size of puzzle map clamped to size of window/drawable area -->
  <svg
    viewBox="0 0 100 100"
    :width="width"
    :height="height"
    ref="el">
  <defs>
    <g id="hex">
      <polygon stroke="#000000" stroke-width="0.5" points="5,-9 -5,-9 -10,0 -5,9 5,9 10,0" />
    </g>
  </defs>
  <g class="hexgrid">
    <use xlink:href="#hex" transform="translate(50, 41)"/>
    <use xlink:href="#hex" transform="translate(35, 50)"/>
    <use xlink:href="#hex" transform="translate(65, 50)"/>
    <use xlink:href="#hex" transform="translate(50, 59)"/>
  </g>
  </svg>
</template>

<style>
/* grid styling */
use {
  transition: 0.4s;
  cursor: pointer;
  fill: transparent;
}
.hexgrid use:hover {
  fill: #0b6942;
}

/* other styling */
body {
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
  margin: 0;
  height: 100vh;
  overflow: hidden;
  font-weight: 700;
  font-family: sans-serif;
}
</style>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { PuzzleJSON, TupleToCoord } from './models'

const el = ref<SVGSVGElement>()
const width = computed(() => window.innerWidth)
const height = computed(() => window.innerHeight)

const griddata : PuzzleJSON = {
  id: 'damm/011',
  name: 'One Way Mirror',
  author: 'Kevin Damm',
  source: 'https://hexoban.com/puzzles/damm/011.json',
  terrain: [
    [-2, -3], [-1, -4], [-1, -3], [-1, -2], [-1, -1],
    [0, -5], [0, -4], [0, -3], [0, -1], [0, 0],
    [1, -3], [1, -2], [1, 0], [1, 1],
    [2, -4], [2, -3], [2, -2], [2, -1], [2, 0], [2, 1], [2, 2],
    [3, -3], [3, -2], [3, 0], [3, 1],
    [4, -2], [4, -1], [4, 1], [4, 2], [4, 3],
    [5, -1], [5, 0], [5, 1], [5, 2], [6, 1]
  ],
  init: {
    goals: [
      [0, -5], [2, -3], [2, -1], [2, 1], [4, -1], [4, 1], [4, 3]
    ],
    crates: [
      [-1, -2], [0, -3], [2, -3], [2, -2], [2, 0], [2, 1], [5, 0]
    ]
  }
}

onMounted(() => {
  for (const tuple of griddata.terrain) {
    console.log('render floor at hex coord ' + TupleToCoord(tuple))
  }
  // viewable grid dimensions
  // const xdim = 4 // Math.floor((width.value) / dx)
  // const ydim = 3 // Math.floor((height.value) / dy)

  // Omitted a brief experiment with doing this on canvas.  I think the interactivity of SVG will come in much more useful here.
})
</script>
