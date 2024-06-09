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

  github:SymbolNotFound/hexoban/webapp/src/components/HexView.vue
-->

<template>
  <!-- TODO should be computing viewBox based on size of puzzle map clamped to size of window/drawable area -->
  <svg
    viewBox="-100 -42 130 144">
  <defs>
    <g id="hex">
      <polygon stroke="#000000" stroke-width="0.5" :points="hexpoints()" />
    </g>
    <g id="goal">
      <circle cx="0" cy="0" r="5.0" stroke-width="1.2" stroke="#D5ADEE" />
    </g>
    <g id="crate">
      <circle cx="0" cy="0" r="3.0" stroke-width="1.2" fill="#D4AF42" />
    </g>
    <g id="at">
      <text x=-5.2 y=0 font-size=".85em"
       style="dominant-baseline: middle"
       stroke-width="0.2" stroke="#000000" fill="#000000">@</text>
    </g>
  </defs>
  <g class="hexgrid">
    <use class="floor" v-for="coord in griddata.terrain" :key="coordID(coord)" xlink:href="#hex" :transform="translate(coord)" />

    <use v-for="goal in griddata.init.goals" :key="goalID(goal)" xlink:href="#goal" :transform="translate(goal)" />
    <use v-for="crate in griddata.init.crates" :key="crateID(crate)" xlink:href="#crate" :transform="translate(crate)" />
    <use xlink:href="#at" :transform="translate(griddata.init.ichiban)" />
  </g>
  </svg>
</template>

<style>
/* grid styling */
use {
  transition: 0.32s;
  cursor: pointer;
  fill: transparent;
}
.hexgrid use.floor:hover {
  fill: #8849d4;
}

/* other styling */
svg { width: 600px; }
body {
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
  margin: 0;
  height: 100vh;
  overflow: hidden;
}
</style>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { PuzzleJSON } from './models'

const el = ref<SVGSVGElement>()
const width = computed(() => window.innerWidth)
const height = computed(() => window.innerHeight)

const w = 9
const h = 15

const hexpoints = () => '-9,-5 -9,5 0,10, 9,5 9,-5 0,-10'
const coordID = (coord: [number, number]) => `hex${coord[0]}_${coord[1]}`
const goalID = (coord: [number, number]) => `goal${coord[0]}_${coord[0]}`
const crateID = (coord: [number, number]) => `crate${coord[0]}_${coord[0]}`

const translate = (val: [number, number]) => {
  const x = (val[1] * 2 * w) - (val[0] * w)
  const y = val[0] * h
  return `translate(${x}, ${y})`
}

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
    ],
    ichiban: [0, 0]
  }
}

onMounted(() => {
  // viewable grid dimensions
  // const xdim = 4 // Math.floor((width.value) / dx)
  // const ydim = 3 // Math.floor((height.value) / dy)
  if (el.value !== undefined) {
    const hv = new HexView(el.value, griddata)
    alert(hv.svg)
  }

  // Omitted a brief experiment with doing this on canvas.  I think the interactivity of SVG will come in much more useful here.
})

// Represents a HexGrid and HexLayout in terms of the viewBox and local
// coordinate transforms of an enclosing SVG element.  This includes any
// panning and zooming done by the client to accommodate viewing a subset
// of the grid as represented by GridLayout.
class HexView {
  // Homogeneous coordinates for applying translation (right column) along
  // with scale, rotation and shear transforms (top-left 1x2 submatrix).  The
  // bottom-right corner element is a normalization factor, rescale it to 0.0.
  transform: [[number, number, number],
              [number, number, number],
              [number, number, number]]

  min: [number, number]
  max: [number, number]

  svg: SVGSVGElement
  viewBox: string

  constructor (el: SVGSVGElement, griddata: PuzzleJSON) {
    // Identity transform (zero translation, uniform projection).
    this.transform = [
      [0, 0, 0],
      [-1, 1, 0],
      [-1, 0, 1]]
    this.svg = el
    if (griddata.terrain.length > 0) {
      this.min = [griddata.terrain[0][0], griddata.terrain[0][1]]
      this.max = [griddata.terrain[0][0], griddata.terrain[0][1]]
    } else {
      this.min = [0, 0]
      this.max = [0, 0]
    }

    this.viewBox = '-42 -51 100 150'
  }
}
</script>
