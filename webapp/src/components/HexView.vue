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
    ref="el"
    :viewBox="(hv && hv.viewBox()) || '0 0 100 100'">
  <defs>
    <!-- floor tiles -->
    <g id="floor">
      <polygon stroke="#000000" stroke-width="0.5" :points="hexpoints()" />
    </g>
    <!-- goal tiles, immovable -->
    <g id="goal">
      <circle cx="0" cy="0" r="5.0" stroke-width="1.2" stroke="#B8B8B8" />
    </g>
    <!-- crates, the movable objects -->
    <g id="crate">
      <!-- <circle cx="0" cy="0" r="3.0" stroke-width="1.2" fill="#FFD700" /> -->
      <text x=-3.4 y=1.2 font-size=".85em"
       dominant-baseline="middle"
       stroke-width="0.2" stroke="#FFD700" fill="#D4AA00">
      $
      </text>
      <text x=-3.4 y=1.2 transform="rotate(90)" font-size=".85em"
       dominant-baseline="middle"
       stroke-width="0.2" stroke="#FFD700" fill="#D4AA00">
      $
      </text>
    </g>
    <!-- the player character -->
    <g id="at">
      <text x=-5.2 y=0 font-size=".85em"
       dominant-baseline="middle"
       stroke-width="0.2" stroke="#D5ADEE" fill="#D5ADEE">
      @
      </text>
    </g>
  </defs>
  <g class="hexgrid">
    <use class="floor" v-for="coord in griddata.terrain" :key="coordID(coord)" xlink:href="#floor" :transform="translate(coord)" />

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
import { ref, onMounted } from 'vue'
import { PuzzleJSON } from './models'

const el = ref<SVGSVGElement>()

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
    ichiban: [0, -1]
  }
}

onMounted(() => {
  if (el.value !== undefined) {
    hv.value = new HexView(el.value, griddata)
    hv.value.Play()
  }

  // viewable grid dimensions
  // const xdim = 4 // Math.floor((window.innerWidth) / dx)
  // const ydim = 3 // Math.floor((window.innerHeight) / dy)
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

  svg: SVGSVGElement
  view: [number, number, number, number]
  viewBox: () => string

  constructor (el: SVGSVGElement, griddata: PuzzleJSON) {
    this.view = [-100, -42, 130, 144]
    // Identity transform (zero translation and identical projection).
    this.transform = [
      [1, 0, 0],
      [0, 1, 0],
      [0, 0, 1]]
    this.svg = el

    if (griddata.terrain.length > 0) {
      const min = [griddata.terrain[0][0], griddata.terrain[0][1]]
      const max = [griddata.terrain[0][0], griddata.terrain[0][1]]

      for (const coord of griddata.terrain) {
        if (coord[0] < min[0]) {
          min[0] = coord[0]
        } else if (coord[0] > max[0]) {
          max[0] = coord[0]
        }

        if (coord[1] < min[1]) {
          min[1] = coord[1]
        } else if (coord[1] > max[1]) {
          max[1] = coord[1]
        }
      }
    }

    this.viewBox = () => `${this.view[0]} ${this.view[1]} ${this.view[2]} ${this.view[3]}`
  }

  Play () {
    addEventListener('replay', (event) => { alert('TODO' + event) })
    addEventListener('stop', (event) => { alert('TODO' + event) })
  }

  Replay () {
    // removeEventListener('replay', this.replayListener )
    addEventListener('play', (event) => { alert('TODO' + event) })
  }

  Stop () {
    addEventListener('stop', (event) => { alert('TODO' + event) })
  }
}

const hv = ref<HexView>()
</script>
