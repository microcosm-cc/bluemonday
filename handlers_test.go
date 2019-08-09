// Copyright (c) 2019, David Kitchen <david@buro9.com></div>
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this
//   list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice,
//   this list of conditions and the following disclaimer in the documentation
//   and/or other materials provided with the distribution.
//
// * Neither the name of the organisation (Microcosm) nor the names of its
//   contributors may be used to endorse or promote products derived from
//   this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package bluemonday

import (
	"sync"
	"testing"
)

func TestDefaultHandlers(t *testing.T) {

	tests := []test{
		{
			in:       `<div style="nonexistentStyle: something;"></div>`,
			expected: `<div></div>`,
		},
		{
			in:       `<div style="aLiGn-cOntEnt: cEntEr;"></div>`,
			expected: `<div style="aLiGn-cOntEnt: cEntEr"></div>`,
		},
		{
			in:       `<div style="align-items: center;"></div>`,
			expected: `<div style="align-items: center"></div>`,
		},
		{
			in:       `<div style="align-self: center;"></div>`,
			expected: `<div style="align-self: center"></div>`,
		},
		{
			in:       `<div style="all: initial;"></div>`,
			expected: `<div style="all: initial"></div>`,
		},
		{
			in: `<div style="animation: mymove 5s infinite;"></div><div` +
				` style="animation: inherit;"></div>`,
			expected: `<div style="animation: mymove 5s infinite"></div><div` +
				` style="animation: inherit"></div>`,
		},
		{
			in: `<div style="animation-delay: 2s;"></div><div style=` +
				`"animation-delay: initial;"></div>`,
			expected: `<div style="animation-delay: 2s"></div><div style=` +
				`"animation-delay: initial"></div>`,
		},
		{
			in:       `<div style="animation-direction: alternate;"></div>`,
			expected: `<div style="animation-direction: alternate"></div>`,
		},
		{
			in: `<div style="animation-duration: 2s;"></div><div style=` +
				`"animation-duration: initial;"></div>`,
			expected: `<div style="animation-duration: 2s"></div><div style=` +
				`"animation-duration: initial"></div>`,
		},
		{
			in:       `<div style="animation-fill-mode: forwards;"></div>`,
			expected: `<div style="animation-fill-mode: forwards"></div>`,
		},
		{
			in: `<div style="animation-iteration-count: 4;"></div><div ` +
				`style="animation-iteration-count: inherit;"></div>`,
			expected: `<div style="animation-iteration-count: 4"></div><div ` +
				`style="animation-iteration-count: inherit"></div>`,
		},
		{
			in: `<div style="animation-name: chuck;"></div><div style=` +
				`"animation-name: none"></div>`,
			expected: `<div style="animation-name: chuck"></div><div style=` +
				`"animation-name: none"></div>`,
		},
		{
			in:       `<div style="animation-play-state: running;"></div>`,
			expected: `<div style="animation-play-state: running"></div>`,
		},
		{
			in: `<div style="animation-timing-function: ` +
				`cubic-bezier(1,1,1,1);"></div><div style=` +
				`"animation-timing-function: steps(2, start);"></div>`,
			expected: `<div style="animation-timing-function: ` +
				`cubic-bezier(1,1,1,1)"></div><div style=` +
				`"animation-timing-function: steps(2, start)"></div>`,
		},
		{
			in:       `<div style="backface-visibility: hidden"></div>`,
			expected: `<div style="backface-visibility: hidden"></div>`,
		},
		{
			in: `<div style="background: lightblue ` +
				`url('https://img_tree.gif') no-repeat fixed center"></div>` +
				`<div style="background: initial"></div>`,
			expected: `<div style="background: lightblue ` +
				`url(&#39;https://img_tree.gif&#39;) no-repeat fixed center">` +
				`</div><div style="background: initial"></div>`,
		},
		{
			in:       `<div style="background-attachment: fixed"></div>`,
			expected: `<div style="background-attachment: fixed"></div>`,
		},
		{
			in:       `<div style="background-blend-mode: lighten"></div>`,
			expected: `<div style="background-blend-mode: lighten"></div>`,
		},
		{
			in:       `<div style="background-clip: padding-box"></div>`,
			expected: `<div style="background-clip: padding-box"></div>`,
		},
		{
			in:       `<div style="background-color: coral"></div>`,
			expected: `<div style="background-color: coral"></div>`,
		},
		{
			in: `<div style="background-image: url('http://paper.gif')">` +
				`</div><div style="background-image: inherit"></div>`,
			expected: `<div style="background-image: ` +
				`url(&#39;http://paper.gif&#39;)"></div><div style="` +
				`background-image: inherit"></div>`,
		},
		{
			in:       `<div style="background-origin: content-box"></div>`,
			expected: `<div style="background-origin: content-box"></div>`,
		},
		{
			in: `<div style="background-position: center"></div><div ` +
				`style="background-position: 20px 20px"></div>`,
			expected: `<div style="background-position: center"></div><div ` +
				`style="background-position: 20px 20px"></div>`,
		},
		{
			in:       `<div style="background-repeat: repeat-y"></div>`,
			expected: `<div style="background-repeat: repeat-y"></div>`,
		},
		{
			in: `<div style="background-size: 300px 100px"></div><div ` +
				`style="background-size: initial"></div>`,
			expected: `<div style="background-size: 300px 100px"></div>` +
				`<div style="background-size: initial"></div>`,
		},
		{
			in: `<div style="border: 4px dotted blue;"></div><div ` +
				`style="border: initial;"></div>`,
			expected: `<div style="border: 4px dotted blue"></div><div ` +
				`style="border: initial"></div>`,
		},
		{
			in: `<div style="border-bottom: 4px dotted blue;"></div>` +
				`<div style="border-bottom: initial"></div>`,
			expected: `<div style="border-bottom: 4px dotted blue"></div>` +
				`<div style="border-bottom: initial"></div>`,
		},
		{
			in:       `<div style="border-bottom-color: blue;"></div>`,
			expected: `<div style="border-bottom-color: blue"></div>`,
		},
		{
			in: `<div style="border-bottom-left-radius: 4px;"></div>` +
				`<div style="border-bottom-left-radius: initial"></div>`,
			expected: `<div style="border-bottom-left-radius: 4px"></div>` +
				`<div style="border-bottom-left-radius: initial"></div>`,
		},
		{
			in: `<div style="border-bottom-right-radius: 40px 4px;">` +
				`</div>`,
			expected: `<div style="border-bottom-right-radius: 40px 4px">` +
				`</div>`,
		},
		{
			in:       `<div style="border-bottom-style: dotted;"></div>`,
			expected: `<div style="border-bottom-style: dotted"></div>`,
		},
		{
			in:       `<div style="border-bottom-width: thin;"></div>`,
			expected: `<div style="border-bottom-width: thin"></div>`,
		},
		{
			in:       `<div style="border-collapse: separate;"></div>`,
			expected: `<div style="border-collapse: separate"></div>`,
		},
		{
			in:       `<div style="border-color: coral;"></div>`,
			expected: `<div style="border-color: coral"></div>`,
		},
		{
			in: `<div style="border-image: url(https://border.png) 30 ` +
				`round;"></div><div style="border-image: initial;"></div>`,
			expected: `<div style="border-image: url(https://border.png) 30 ` +
				`round"></div><div style="border-image: initial"></div>`,
		},
		{
			in:       `<div style="border-image-outset: 10px;"></div>`,
			expected: `<div style="border-image-outset: 10px"></div>`,
		},
		{
			in:       `<div style="border-image-repeat: repeat;"></div>`,
			expected: `<div style="border-image-repeat: repeat"></div>`,
		},
		{
			in: `<div style="border-image-slice: 30%;"></div><div ` +
				`style="border-image-slice: fill;"></div><div style="` +
				`border-image-slice: 3% 3% 3% 3% 3%;"></div>`,
			expected: `<div style="border-image-slice: 30%"></div><div style` +
				`="border-image-slice: fill"></div><div></div>`,
		},
		{
			in: `<div style="border-image-source: ` +
				`url(https://border.png);"></div>`,
			expected: `<div style="border-image-source: ` +
				`url(https://border.png)"></div>`,
		},
		{
			in:       `<div style="border-image-width: 10px;"></div>`,
			expected: `<div style="border-image-width: 10px"></div>`,
		},
		{
			in:       `<div style="border-left: 4px dotted blue;"></div>`,
			expected: `<div style="border-left: 4px dotted blue"></div>`,
		},
		{
			in:       `<div style="border-left-color: blue;"></div>`,
			expected: `<div style="border-left-color: blue"></div>`,
		},
		{
			in:       `<div style="border-left-style: dotted;"></div>`,
			expected: `<div style="border-left-style: dotted"></div>`,
		},
		{
			in:       `<div style="border-left-width: thin;"></div>`,
			expected: `<div style="border-left-width: thin"></div>`,
		},
		{
			in: `<div style="border-radius: 25px;"></div><div style=` +
				`"border-radius: initial;"></div><div style="border-radius:` +
				` 1px 1px 1px 1px 1px;"></div>`,
			expected: `<div style="border-radius: 25px"></div><div style=` +
				`"border-radius: initial"></div><div></div>`,
		},
		{
			in:       `<div style="border-left: 4px dotted blue;"></div>`,
			expected: `<div style="border-left: 4px dotted blue"></div>`,
		},
		{
			in:       `<div style="border-right-color: blue;"></div>`,
			expected: `<div style="border-right-color: blue"></div>`,
		},
		{
			in:       `<div style="border-right-style: dotted;"></div>`,
			expected: `<div style="border-right-style: dotted"></div>`,
		},
		{
			in:       `<div style="border-right-width: thin;"></div>`,
			expected: `<div style="border-right-width: thin"></div>`,
		},
		{
			in:       `<div style="border-spacing: 15px;"></div>`,
			expected: `<div style="border-spacing: 15px"></div>`,
		},
		{
			in: `<div style="border-style: dotted;"></div><div style="` +
				`border-style: initial;"></div><div style="border-style: ` +
				`dotted dotted dotted dotted dotted;"></div>`,
			expected: `<div style="border-style: dotted"></div><div style=` +
				`"border-style: initial"></div><div></div>`,
		},
		{
			in:       `<div style="border-top: 4px dotted blue;"></div>`,
			expected: `<div style="border-top: 4px dotted blue"></div>`,
		},
		{
			in:       `<div style="border-top-color: blue;"></div>`,
			expected: `<div style="border-top-color: blue"></div>`,
		},
		{
			in:       `<div style="border-top-left-radius: 4px;"></div>`,
			expected: `<div style="border-top-left-radius: 4px"></div>`,
		},
		{
			in:       `<div style="border-top-right-radius: 40px 4px;"></div>`,
			expected: `<div style="border-top-right-radius: 40px 4px"></div>`,
		},
		{
			in:       `<div style="border-top-style: dotted;"></div>`,
			expected: `<div style="border-top-style: dotted"></div>`,
		},
		{
			in:       `<div style="border-top-width: thin;"></div>`,
			expected: `<div style="border-top-width: thin"></div>`,
		},
		{
			in: `<div style="border-width: thin;"></div><div style="` +
				`border-width: initial;"></div><div style="border-width: ` +
				`thin thin thin thin thin;"></div>`,
			expected: `<div style="border-width: thin"></div><div style="` +
				`border-width: initial"></div><div></div>`,
		},
		{
			in: `<div style="bottom: 10px;"></div><div style="bottom:` +
				` auto;"></div>`,
			expected: `<div style="bottom: 10px"></div><div style="bottom:` +
				` auto"></div>`,
		},
		{
			in:       `<div style="box-decoration-break: slice;"></div>`,
			expected: `<div style="box-decoration-break: slice"></div>`,
		},
		{
			in: `<div style="box-shadow: 10px 10px #888888;"></div>` +
				`<div style="box-shadow: aa;"></div><div style="box-shadow: ` +
				`10px aa;"></div><div style="box-shadow: 10px;"></div><div ` +
				`style="box-shadow: 10px 10px aa;"></div>`,
			expected: `<div style="box-shadow: 10px 10px #888888"></div>` +
				`<div></div><div></div><div></div><div></div>`,
		},
		{
			in:       `<div style="box-sizing: border-box;"></div>`,
			expected: `<div style="box-sizing: border-box"></div>`,
		},
		{
			in:       `<div style="break-after: column;"></div>`,
			expected: `<div style="break-after: column"></div>`,
		},
		{
			in:       `<div style="break-before: column;"></div>`,
			expected: `<div style="break-before: column"></div>`,
		},
		{
			in:       `<div style="break-inside: avoid-column;"></div>`,
			expected: `<div style="break-inside: avoid-column"></div>`,
		},
		{
			in:       `<div style="caption-side: bottom;"></div>`,
			expected: `<div style="caption-side: bottom"></div>`,
		},
		{
			in: `<div style="caret-color: red;"></div><div style=` +
				`"caret-color: rgb(2,2,2);"></div><div style="caret-color:` +
				` rgba(2,2,2,0.5);"></div><div style="caret-color: ` +
				`hsl(2,2%,2%);"></div><div style="caret-color: ` +
				`hsla(2,2%,2%,0.5);"></div>`,
			expected: `<div style="caret-color: red"></div><div style=` +
				`"caret-color: rgb(2,2,2)"></div><div style="caret-color: ` +
				`rgba(2,2,2,0.5)"></div><div style="caret-color: ` +
				`hsl(2,2%,2%)"></div><div style="caret-color: ` +
				`hsla(2,2%,2%,0.5)"></div>`,
		},
		{
			in:       `<div style="clear: both;"></div>`,
			expected: `<div style="clear: both"></div>`,
		},
		{
			in: `<div style="clip: rect(0px,60px,200px,0px);"></div>` +
				`<div style="clip: auto;"></div>`,
			expected: `<div style="clip: rect(0px,60px,200px,0px)"></div>` +
				`<div style="clip: auto"></div>`,
		},
		{
			in: `<div style="color: red;"></div><div style="color: ` +
				`rgb(2,2,2);"></div><div style="color: rgba(2,2,2,0.5);">` +
				`</div><div style="color: hsl(2,2%,2%);"></div><div style="` +
				`color: hsla(2,2%,2%,0.5);"></div>`,
			expected: `<div style="color: red"></div><div style="color: ` +
				`rgb(2,2,2)"></div><div style="color: rgba(2,2,2,0.5)">` +
				`</div><div style="color: hsl(2,2%,2%)"></div><div style="` +
				`color: hsla(2,2%,2%,0.5)"></div>`,
		},
		{
			in:       `<div style="clear: both;"></div>`,
			expected: `<div style="clear: both"></div>`,
		},
		{
			in: `<div style="column-count: 3;"></div><div style="` +
				`column-count: auto;"></div>`,
			expected: `<div style="column-count: 3"></div><div style="` +
				`column-count: auto"></div>`,
		},
		{
			in:       `<div style="column-fill: balance;"></div>`,
			expected: `<div style="column-fill: balance"></div>`,
		},
		{
			in: `<div style="column-gap: 40px;"></div><div style="` +
				`column-gap: normal;"></div>`,
			expected: `<div style="column-gap: 40px"></div><div style="` +
				`column-gap: normal"></div>`,
		},
		{
			in:       `<div style="column-rule: 4px double #ff00ff;"></div>`,
			expected: `<div style="column-rule: 4px double #ff00ff"></div>`,
		},
		{
			in:       `<div style="column-rule-color: #ff00ff;"></div>`,
			expected: `<div style="column-rule-color: #ff00ff"></div>`,
		},
		{
			in:       `<div style="column-rule: red;"></div>`,
			expected: `<div style="column-rule: red"></div>`,
		},
		{
			in:       `<div style="column-rule-width: 4px;"></div>`,
			expected: `<div style="column-rule-width: 4px"></div>`,
		},
		{
			in:       `<div style="column-span: all;"></div>`,
			expected: `<div style="column-span: all"></div>`,
		},
		{
			in: `<div style="column-width: 4px;"></div><div style="` +
				`column-width: auto;"></div>`,
			expected: `<div style="column-width: 4px"></div><div style="` +
				`column-width: auto"></div>`,
		},
		{
			in: `<div style="columns: 4px 3"></div><div style="` +
				`columns: auto"></div>`,
			expected: `<div style="columns: 4px 3"></div><div style="` +
				`columns: auto"></div>`,
		},
		{
			in:       `<div style="cursor: alias"></div>`,
			expected: `<div style="cursor: alias"></div>`,
		},
		{
			in:       `<div style="direction: rtl"></div>`,
			expected: `<div style="direction: rtl"></div>`,
		},
		{
			in:       `<div style="display: block"></div>`,
			expected: `<div style="display: block"></div>`,
		},
		{
			in:       `<div style="empty-cells: hide"></div>`,
			expected: `<div style="empty-cells: hide"></div>`,
		},
		{
			in: `<div style="filter: grayscale(100%)"></div><div style` +
				`="filter: sepia(100%)"></div>`,
			expected: `<div style="filter: grayscale(100%)"></div><div style` +
				`="filter: sepia(100%)"></div>`,
		},
		{
			in: `<div style="flex: 1"></div><div style="flex: auto">` +
				`</div>`,
			expected: `<div style="flex: 1"></div><div style="flex: auto">` +
				`</div>`,
		},
		{
			in: `<div style="flex-basis: 10px"></div><div style="` +
				`flex-basis: auto"></div>`,
			expected: `<div style="flex-basis: 10px"></div><div style="` +
				`flex-basis: auto"></div>`,
		},
		{
			in:       `<div style="flex-direction: row-reverse"></div>`,
			expected: `<div style="flex-direction: row-reverse"></div>`,
		},
		{
			in: `<div style="flex-flow: row-reverse wrap"></div><div ` +
				`style="flex-flow: initial"></div>`,
			expected: `<div style="flex-flow: row-reverse wrap"></div><div ` +
				`style="flex-flow: initial"></div>`,
		},
		{
			in: `<div style="flex-grow: 1"></div><div style="flex-grow` +
				`: initial"></div>`,
			expected: `<div style="flex-grow: 1"></div><div style="flex-grow` +
				`: initial"></div>`,
		},
		{
			in:       `<div style="flex-shrink: 3"></div>`,
			expected: `<div style="flex-shrink: 3"></div>`,
		},
		{
			in:       `<div style="flex-wrap: wrap"></div>`,
			expected: `<div style="flex-wrap: wrap"></div>`,
		},
		{
			in:       `<div style="float: right"></div>`,
			expected: `<div style="float: right"></div>`,
		},
		{
			in: `<div style="font: italic bold 12px/30px Georgia, serif` +
				`"></div><div style="font: icon"></div>`,
			expected: `<div style="font: italic bold 12px/30px Georgia, serif` +
				`"></div><div style="font: icon"></div>`,
		},
		{
			in: `<div style="font-family: 'Times New Roman', Times, ` +
				`serif"></div><span style="font-family: comic sans ms, ` +
				`cursive, sans-serif;">aaaaaa</span></span>`,
			expected: `<div style="font-family: &#39;Times New Roman&#39;,` +
				` Times, serif"></div><span style="font-family: comic sans` +
				` ms, cursive, sans-serif">aaaaaa</span></span>`,
		},
		{
			in:       `<div style="font-kerning: normal"></div>`,
			expected: `<div style="font-kerning: normal"></div>`,
		},
		{
			in:       `<div style="font-language-override: normal"></div>`,
			expected: `<div style="font-language-override: normal"></div>`,
		},
		{
			in:       `<div style="font-size: large"></div>`,
			expected: `<div style="font-size: large"></div>`,
		},
		{
			in: `<div style="font-size-adjust: 0.58"></div><div style="` +
				`font-size-adjust: auto"></div>`,
			expected: `<div style="font-size-adjust: 0.58"></div><div style="` +
				`font-size-adjust: auto"></div>`,
		},
		{
			in:       `<div style="font-stretch: expanded"></div>`,
			expected: `<div style="font-stretch: expanded"></div>`,
		},
		{
			in:       `<div style="font-style: italic"></div>`,
			expected: `<div style="font-style: italic"></div>`,
		},
		{
			in:       `<div style="font-synthesis: style"></div>`,
			expected: `<div style="font-synthesis: style"></div>`,
		},
		{
			in:       `<div style="font-variant: small-caps"></div>`,
			expected: `<div style="font-variant: small-caps"></div>`,
		},
		{
			in:       `<div style="font-variant-caps: small-caps"></div>`,
			expected: `<div style="font-variant-caps: small-caps"></div>`,
		},
		{
			in:       `<div style="font-variant-position: sub"></div>`,
			expected: `<div style="font-variant-position: sub"></div>`,
		},
		{
			in:       `<div style="font-weight: normal"></div>`,
			expected: `<div style="font-weight: normal"></div>`,
		},
		{
			in: `<div style="grid: 150px / auto auto auto;"></div><div ` +
				`style="grid: none;"></div>`,
			expected: `<div style="grid: 150px / auto auto auto"></div><div ` +
				`style="grid: none"></div>`,
		},
		{
			in: `<div style="grid-area: 2 / 1 / span 2 / span 3;">` +
				`</div>`,
			expected: `<div style="grid-area: 2 / 1 / span 2 / span 3">` +
				`</div>`,
		},
		{
			in: `<div style="grid-auto-columns: 150px;"></div>` +
				`<div style="grid-auto-columns: auto;"></div>`,
			expected: `<div style="grid-auto-columns: 150px"></div>` +
				`<div style="grid-auto-columns: auto"></div>`,
		},
		{
			in:       `<div style="grid-auto-flow: column;"></div>`,
			expected: `<div style="grid-auto-flow: column"></div>`,
		},
		{
			in:       `<div style="grid-auto-rows: 150px;"></div>`,
			expected: `<div style="grid-auto-rows: 150px"></div>`,
		},
		{
			in:       `<div style="grid-column: 1 / span 2;"></div>`,
			expected: `<div style="grid-column: 1 / span 2"></div>`,
		},
		{
			in: `<div style="grid-column-end: span 2;"></div>` +
				`<div style="grid-column-end: auto;"></div>`,
			expected: `<div style="grid-column-end: span 2"></div>` +
				`<div style="grid-column-end: auto"></div>`,
		},
		{
			in:       `<div style="grid-column-gap: 10px;"></div>`,
			expected: `<div style="grid-column-gap: 10px"></div>`,
		},
		{
			in:       `<div style="grid-column-start: 1;"></div>`,
			expected: `<div style="grid-column-start: 1"></div>`,
		},
		{
			in: `<div style="grid-gap: 1px;"></div><div style=` +
				`"grid-gap: 1px 1px 1px;"></div>`,
			expected: `<div style="grid-gap: 1px"></div><div></div>`,
		},
		{
			in:       `<div style="grid-row: 1 / span 2;"></div>`,
			expected: `<div style="grid-row: 1 / span 2"></div>`,
		},
		{
			in:       `<div style="grid-row-end: span 2;"></div>`,
			expected: `<div style="grid-row-end: span 2"></div>`,
		},
		{
			in:       `<div style="grid-row-gap: 10px;"></div>`,
			expected: `<div style="grid-row-gap: 10px"></div>`,
		},
		{
			in:       `<div style="grid-row-start: 1;"></div>`,
			expected: `<div style="grid-row-start: 1"></div>`,
		},
		{
			in: `<div style="grid-template: 150px / auto auto auto;">` +
				`</div><div style="grid-template: none"></div><div style="` +
				`grid-template: a / a / a"></div>`,
			expected: `<div style="grid-template: 150px / auto auto auto">` +
				`</div><div style="grid-template: none"></div><div></div>`,
		},
		{
			in: `<div style="grid-template-areas: none;"></div><div ` +
				`style="grid-template-areas: 'Billy'"></div>`,
			expected: `<div style="grid-template-areas: none"></div>` +
				`<div style="grid-template-areas: &#39;Billy&#39;"></div>`,
		},
		{
			in: `<div style="grid-template-columns: auto auto auto` +
				` auto auto;"></div>`,
			expected: `<div style="grid-template-columns: auto auto` +
				` auto auto auto"></div>`,
		},
		{
			in: `<div style="grid-template-rows: 150px 150px">` +
				`</div><div style="grid-template-rows: aaaa aaaaa"></div>`,
			expected: `<div style="grid-template-rows: 150px 150px">` +
				`</div><div></div>`,
		},
		{
			in:       `<div style="hanging-punctuation: first;"></div>`,
			expected: `<div style="hanging-punctuation: first"></div>`,
		},
		{
			in: `<div style="height: 50px;"></div><div style="height: ` +
				`auto;"></div>`,
			expected: `<div style="height: 50px"></div><div style="height: ` +
				`auto"></div>`,
		},
		{
			in:       `<div style="hyphens: manual;"></div>`,
			expected: `<div style="hyphens: manual"></div>`,
		},
		{
			in:       `<div style="isolation: isolate;"></div>`,
			expected: `<div style="isolation: isolate"></div>`,
		},
		{
			in:       `<div style="image-rendering: smooth;"></div>`,
			expected: `<div style="image-rendering: smooth"></div>`,
		},
		{
			in:       `<div style="justify-content: center;"></div>`,
			expected: `<div style="justify-content: center"></div>`,
		},
		{
			in:       `<div style="left: 150px;"></div>`,
			expected: `<div style="left: 150px"></div>`,
		},
		{
			in: `<div style="letter-spacing: -3px;"></div><div style` +
				`="letter-spacing: normal;"></div>`,
			expected: `<div style="letter-spacing: -3px"></div><div style` +
				`="letter-spacing: normal"></div>`,
		},
		{
			in:       `<div style="line-break: auto"></div>`,
			expected: `<div style="line-break: auto"></div>`,
		},
		{
			in: `<div style="line-height: 1.6;"></div><div style=` +
				`"line-height: normal;"></div>`,
			expected: `<div style="line-height: 1.6"></div><div style=` +
				`"line-height: normal"></div>`,
		},
		{
			in: `<div style="list-style: square inside ` +
				`url(http://sqpurple.gif);"></div><div style="list-style: ` +
				`initial"></div>`,
			expected: `<div style="list-style: square inside ` +
				`url(http://sqpurple.gif)"></div><div style="list-style: ` +
				`initial"></div>`,
		},
		{
			in: `<div style="list-style-image: ` +
				`url(http://sqpurple.gif);"></div>`,
			expected: `<div style="list-style-image: ` +
				`url(http://sqpurple.gif)"></div>`,
		},
		{
			in:       `<div style="list-style-position: inside;"></div>`,
			expected: `<div style="list-style-position: inside"></div>`,
		},
		{
			in:       `<div style="list-style-type: square;"></div>`,
			expected: `<div style="list-style-type: square"></div>`,
		},
		{
			in: `<div style="margin: 150px;"></div><div style="margin:` +
				` auto;"></div>`,
			expected: `<div style="margin: 150px"></div><div style="margin:` +
				` auto"></div>`,
		},
		{
			in: `<div style="margin-bottom: 150px;"></div><div ` +
				`style="margin-bottom: auto;"></div>`,
			expected: `<div style="margin-bottom: 150px"></div><div ` +
				`style="margin-bottom: auto"></div>`,
		},
		{
			in:       `<div style="margin-left: 150px;"></div>`,
			expected: `<div style="margin-left: 150px"></div>`,
		},
		{
			in:       `<div style="margin-right: 150px;"></div>`,
			expected: `<div style="margin-right: 150px"></div>`,
		},
		{
			in:       `<div style="margin-top: 150px;"></div>`,
			expected: `<div style="margin-top: 150px"></div>`,
		},
		{
			in: `<div style="max-height: 150px;"></div><div style=` +
				`"max-height: initial;"></div>`,
			expected: `<div style="max-height: 150px"></div><div style=` +
				`"max-height: initial"></div>`,
		},
		{
			in:       `<div style="max-width: 150px;"></div>`,
			expected: `<div style="max-width: 150px"></div>`,
		},
		{
			in: `<div style="min-height: 150px;"></div><div style=` +
				`"min-height: initial;"></div>`,
			expected: `<div style="min-height: 150px"></div><div style=` +
				`"min-height: initial"></div>`,
		},
		{
			in:       `<div style="min-width: 150px;"></div>`,
			expected: `<div style="min-width: 150px"></div>`,
		},
		{
			in:       `<div style="mix-blend-mode: darken;"></div>`,
			expected: `<div style="mix-blend-mode: darken"></div>`,
		},
		{
			in:       `<div style="object-fit: cover;"></div>`,
			expected: `<div style="object-fit: cover"></div>`,
		},
		{
			in: `<div style="object-position: 5px 10%;"></div><div ` +
				`style="object-position: initial"></div><div style="` +
				`object-position: 5px 10% 5px;"></div>`,
			expected: `<div style="object-position: 5px 10%"></div><div ` +
				`style="object-position: initial"></div><div></div>`,
		},
		{
			in: `<div style="opacity: 0.5;"></div><div style="opacity:` +
				` initial"></div>`,
			expected: `<div style="opacity: 0.5"></div><div style="opacity:` +
				` initial"></div>`,
		},
		{
			in: `<div style="order: 2;"></div><div style="order: ` +
				`initial"></div>`,
			expected: `<div style="order: 2"></div><div style="order: ` +
				`initial"></div>`,
		},
		{
			in: `<div style="outline: 2px dashed blue;"></div><div ` +
				`style="outline: initial"></div>`,
			expected: `<div style="outline: 2px dashed blue"></div><div ` +
				`style="outline: initial"></div>`,
		},
		{
			in:       `<div style="outline-color: blue;"></div>`,
			expected: `<div style="outline-color: blue"></div>`,
		},
		{
			in: `<div style="outline-offset: 2px;"></div><div ` +
				`style="outline-offset: initial;"></div>`,
			expected: `<div style="outline-offset: 2px"></div><div ` +
				`style="outline-offset: initial"></div>`,
		},
		{
			in:       `<div style="outline-style: dashed;"></div>`,
			expected: `<div style="outline-style: dashed"></div>`,
		},
		{
			in:       `<div style="outline-width: thick;"></div>`,
			expected: `<div style="outline-width: thick"></div>`,
		},
		{
			in:       `<div style="overflow: scroll;"></div>`,
			expected: `<div style="overflow: scroll"></div>`,
		},
		{
			in:       `<div style="overflow-x: scroll;"></div>`,
			expected: `<div style="overflow-x: scroll"></div>`,
		},
		{
			in:       `<div style="overflow-y: scroll;"></div>`,
			expected: `<div style="overflow-y: scroll"></div>`,
		},
		{
			in:       `<div style="overflow-wrap: anywhere;"></div>`,
			expected: `<div style="overflow-wrap: anywhere"></div>`,
		},
		{
			in:       `<div style="orphans: 2;"></div>`,
			expected: `<div style="orphans: 2"></div>`,
		},
		{
			in:       `<div style="padding: 55px;"></div>`,
			expected: `<div style="padding: 55px"></div>`,
		},
		{
			in: `<div style="padding-bottom: 55px;"></div><div style` +
				`="padding-bottom: initial;"></div>`,
			expected: `<div style="padding-bottom: 55px"></div><div style=` +
				`"padding-bottom: initial"></div>`,
		},
		{
			in:       `<div style="padding-left: 55px;"></div>`,
			expected: `<div style="padding-left: 55px"></div>`,
		},
		{
			in:       `<div style="padding-right: 55px;"></div>`,
			expected: `<div style="padding-right: 55px"></div>`,
		},
		{
			in:       `<div style="padding-top: 55px;"></div>`,
			expected: `<div style="padding-top: 55px"></div>`,
		},
		{
			in:       `<div style="page-break-after: always;"></div>`,
			expected: `<div style="page-break-after: always"></div>`,
		},
		{
			in:       `<div style="page-break-before: always;"></div>`,
			expected: `<div style="page-break-before: always"></div>`,
		},
		{
			in:       `<div style="page-break-inside: avoid;"></div>`,
			expected: `<div style="page-break-inside: avoid"></div>`,
		},
		{
			in: `<div style="perspective: 100px;"></div><div style=` +
				`"perspective: none;"></div>`,
			expected: `<div style="perspective: 100px"></div><div style=` +
				`"perspective: none"></div>`,
		},
		{
			in:       `<div style="perspective-origin: left;"></div>`,
			expected: `<div style="perspective-origin: left"></div>`,
		},
		{
			in:       `<div style="pointer-events: auto;"></div>`,
			expected: `<div style="pointer-events: auto"></div>`,
		},
		{
			in:       `<div style="position: absolute;"></div>`,
			expected: `<div style="position: absolute"></div>`,
		},
		{
			in:       `<div style="quotes: '‹' '›';"></div>`,
			expected: `<div style="quotes: &#39;‹&#39; &#39;›&#39;"></div>`,
		},
		{
			in:       `<div style="resize: both;"></div>`,
			expected: `<div style="resize: both"></div>`,
		},
		{
			in:       `<div style="right: 10px;"></div>`,
			expected: `<div style="right: 10px"></div>`,
		},
		{
			in:       `<div style="scroll-behavior: smooth;"></div>`,
			expected: `<div style="scroll-behavior: smooth"></div>`,
		},
		{
			in: `<div style="tab-size: 16;"></div><div style="tab-size:` +
				` initial;"></div>`,
			expected: `<div style="tab-size: 16"></div><div style="tab-size:` +
				` initial"></div>`,
		},
		{
			in:       `<div style="table-layout: fixed;"></div>`,
			expected: `<div style="table-layout: fixed"></div>`,
		},
		{
			in:       `<div style="text-align: justify;"></div>`,
			expected: `<div style="text-align: justify"></div>`,
		},
		{
			in:       `<div style="text-align-last: justify;"></div>`,
			expected: `<div style="text-align-last: justify"></div>`,
		},
		{
			in: `<div style="text-combine-upright: none;"></div><div` +
				` style="text-combine-upright: digits 2"></div>`,
			expected: `<div style="text-combine-upright: none"></div><div ` +
				`style="text-combine-upright: digits 2"></div>`,
		},
		{
			in: `<div style="text-decoration: underline underline;">` +
				`</div><div style="text-decoration: initial"></div>`,
			expected: `<div style="text-decoration: underline underline">` +
				`</div><div style="text-decoration: initial"></div>`,
		},
		{
			in:       `<div style="text-decoration-color: red;"></div>`,
			expected: `<div style="text-decoration-color: red"></div>`,
		},
		{
			in: `<div style="text-decoration-line: underline ` +
				`underline;"></div>`,
			expected: `<div style="text-decoration-line: underline ` +
				`underline"></div>`,
		},
		{
			in:       `<div style="text-decoration-style: solid;"></div>`,
			expected: `<div style="text-decoration-style: solid"></div>`,
		},
		{
			in: `<div style="text-indent: 30%;"></div><div style=` +
				`"text-indent: initial"></div>`,
			expected: `<div style="text-indent: 30%"></div><div style=` +
				`"text-indent: initial"></div>`,
		},
		{
			in:       `<div style="text-orientation: mixed"></div>`,
			expected: `<div style="text-orientation: mixed"></div>`,
		},
		{
			in:       `<div style="text-justify: inter-word;"></div>`,
			expected: `<div style="text-justify: inter-word"></div>`,
		},
		{
			in: `<div style="text-overflow: ellipsis;"></div><div ` +
				`style="text-overflow: 'something'"></div>`,
			expected: `<div style="text-overflow: ellipsis"></div><div ` +
				`style="text-overflow: &#39;something&#39;"></div>`,
		},
		{
			in:       `<div style="text-shadow: 2px 2px #ff0000;"></div>`,
			expected: `<div style="text-shadow: 2px 2px #ff0000"></div>`,
		},
		{
			in:       `<div style="text-transform: uppercase;"></div>`,
			expected: `<div style="text-transform: uppercase"></div>`,
		},
		{
			in:       `<div style="top: 150px;"></div>`,
			expected: `<div style="top: 150px"></div>`,
		},
		{
			in: `<div style="transform: scaleY(1.5);"></div><div ` +
				`style="transform: perspective(20px);"></div>`,
			expected: `<div style="transform: scaleY(1.5)"></div><div ` +
				`style="transform: perspective(20px)"></div>`,
		},
		{
			in:       `<div style="transform-origin: 40% 40%;"></div>`,
			expected: `<div style="transform-origin: 40% 40%"></div>`,
		},
		{
			in:       `<div style="transform-style: preserve-3d;"></div>`,
			expected: `<div style="transform-style: preserve-3d"></div>`,
		},
		{
			in:       `<div style="transition: width 2s;"></div>`,
			expected: `<div style="transition: width 2s"></div>`,
		},
		{
			in: `<div style="transition-delay: 2s;"></div><div ` +
				`style="transition-delay: initial;"></div>`,
			expected: `<div style="transition-delay: 2s"></div><div ` +
				`style="transition-delay: initial"></div>`,
		},
		{
			in: `<div style="transition-duration: 2s;"></div><div ` +
				`style="transition-duration: initial;"></div>`,
			expected: `<div style="transition-duration: 2s"></div><div ` +
				`style="transition-duration: initial"></div>`,
		},
		{
			in: `<div style="transition-property: width;"></div><div ` +
				`style="transition-property: initial;"></div>`,
			expected: `<div style="transition-property: width"></div><div ` +
				`style="transition-property: initial"></div>`,
		},
		{
			in: `<div style="transition-timing-function: linear;">` +
				`</div>`,
			expected: `<div style="transition-timing-function: linear">` +
				`</div>`,
		},
		{
			in:       `<div style="unicode-bidi: bidi-override;"></div>`,
			expected: `<div style="unicode-bidi: bidi-override"></div>`,
		},
		{
			in:       `<div style="user-select: none;"></div>`,
			expected: `<div style="user-select: none"></div>`,
		},
		{
			in:       `<div style="vertical-align: text-bottom;"></div>`,
			expected: `<div style="vertical-align: text-bottom"></div>`,
		},
		{
			in:       `<div style="visibility: visible;"></div>`,
			expected: `<div style="visibility: visible"></div>`,
		},
		{
			in:       `<div style="white-space: normal;"></div>`,
			expected: `<div style="white-space: normal"></div>`,
		},
		{
			in: `<div style="width: 130px;"></div><div style="width: ` +
				`auto;"></div>`,
			expected: `<div style="width: 130px"></div><div style="width: ` +
				`auto"></div>`,
		},
		{
			in:       `<div style="word-break: break-all;"></div>`,
			expected: `<div style="word-break: break-all"></div>`,
		},
		{
			in: `<div style="word-spacing: 30px;"></div><div style=` +
				`"word-spacing: normal"></div>`,
			expected: `<div style="word-spacing: 30px"></div><div style=` +
				`"word-spacing: normal"></div>`,
		},
		{
			in:       `<div style="word-wrap: break-word;"></div>`,
			expected: `<div style="word-wrap: break-word"></div>`,
		},
		{
			in:       `<div style="writing-mode: vertical-rl;"></div>`,
			expected: `<div style="writing-mode: vertical-rl"></div>`,
		},
		{
			in: `<div style="z-index: -1;"></div><div style="z-index:` +
				` auto;"></div>`,
			expected: `<div style="z-index: -1"></div><div style="z-index:` +
				` auto"></div>`,
		},
	}

	p := UGCPolicy()
	p.AllowStyles("nonexistentStyle", "align-content", "align-items",
		"align-self", "all", "animation", "animation-delay",
		"animation-direction", "animation-duration", "animation-fill-mode",
		"animation-iteration-count", "animation-name", "animation-play-state",
		"animation-timing-function", "backface-visibility", "background",
		"background-attachment", "background-blend-mode", "background-clip",
		"background-color", "background-image", "background-origin",
		"background-position", "background-repeat", "background-size",
		"border", "border-bottom", "border-bottom-color",
		"border-bottom-left-radius", "border-bottom-right-radius",
		"border-bottom-style", "border-bottom-width", "border-collapse",
		"border-color", "border-image", "border-image-outset",
		"border-image-repeat", "border-image-slice", "border-image-source",
		"border-image-width", "border-left", "border-left-color",
		"border-left-style", "border-left-width", "border-radius",
		"border-right", "border-right-color", "border-right-style",
		"border-right-width", "border-spacing", "border-style", "border-top",
		"border-top-color", "border-top-left-radius",
		"border-top-right-radius", "border-top-style", "border-top-width",
		"border-width", "bottom", "box-decoration-break", "box-shadow",
		"box-sizing", "break-after", "break-before", "break-inside",
		"caption-side", "caret-color", "clear", "clip", "color",
		"column-count", "column-fill", "column-gap", "column-rule",
		"column-rule-color", "column-rule-style", "column-rule-width",
		"column-span", "column-width", "columns", "cursor", "direction",
		"display", "empty-cells", "filter", "flex", "flex-basis",
		"flex-direction", "flex-flow", "flex-grow", "flex-shrink",
		"flex-wrap", "float", "font", "font-family", "font-kerning",
		"font-language-override", "font-size", "font-size-adjust",
		"font-stretch", "font-style", "font-synthesis", "font-variant",
		"font-variant-caps", "font-variant-position", "font-weight", "grid",
		"grid-area", "grid-auto-columns", "grid-auto-flow", "grid-auto-rows",
		"grid-column", "grid-column-end", "grid-column-gap",
		"grid-column-start", "grid-gap", "grid-row", "grid-row-end",
		"grid-row-gap", "grid-row-start", "grid-template",
		"grid-template-areas", "grid-template-columns", "grid-template-rows",
		"hanging-punctuation", "height", "hyphens", "image-rendering",
		"isolation", "justify-content", "left", "letter-spacing", "line-break",
		"line-height", "list-style", "list-style-image", "list-style-position",
		"list-style-type", "margin", "margin-bottom", "margin-left",
		"margin-right", "margin-top", "max-height", "max-width", "min-height",
		"min-width", "mix-blend-mode", "object-fit", "object-position",
		"opacity", "order", "orphans", "outline", "outline-color",
		"outline-offset", "outline-style", "outline-width", "overflow",
		"overflow-wrap", "overflow-x", "overflow-y", "padding",
		"padding-bottom", "padding-left", "padding-right", "padding-top",
		"page-break-after", "page-break-before", "page-break-inside",
		"perspective", "perspective-origin", "pointer-events", "position",
		"quotes", "resize", "right", "scroll-behavior", "tab-size",
		"table-layout", "text-align", "text-align-last",
		"text-combine-upright", "text-decoration", "text-decoration-color",
		"text-decoration-line", "text-decoration-style", "text-indent",
		"text-justify", "text-orientation", "text-overflow", "text-shadow",
		"text-transform", "top", "transform", "transform-origin",
		"transform-style", "transition", "transition-delay",
		"transition-duration", "transition-property",
		"transition-timing-function", "unicode-bidi", "user-select",
		"vertical-align", "visibility", "white-space", "widows", "width",
		"word-break", "word-spacing", "word-wrap", "writing-mode",
		"z-index").Globally()
	p.RequireParseableURLs(true)

	// These tests are run concurrently to enable the race detector to pick up
	// potential issues
	wg := sync.WaitGroup{}
	wg.Add(len(tests))
	for ii, tt := range tests {
		go func(ii int, tt test) {
			out := p.Sanitize(tt.in)
			if out != tt.expected {
				t.Errorf(
					"test %d failed;\ninput   : %s\noutput  : %s\nexpected: %s",
					ii,
					tt.in,
					out,
					tt.expected,
				)
			}
			wg.Done()
		}(ii, tt)
	}
	wg.Wait()
}

func TestUnicodePoints(t *testing.T) {

	tests := []test{
		{
			in:       `<div style="color: \72 ed;"></div>`,
			expected: `<div style="color: \72 ed"></div>`,
		},
		{
			in:       `<div style="color: \0072 ed;"></div>`,
			expected: `<div style="color: \0072 ed"></div>`,
		},
		{
			in:       `<div style="color: \000072 ed;"></div>`,
			expected: `<div style="color: \000072 ed"></div>`,
		},
		{
			in:       `<div style="color: \000072ed;"></div>`,
			expected: `<div style="color: \000072ed"></div>`,
		},
		{
			in:       `<div style="color: \100072ed;"></div>`,
			expected: `<div></div>`,
		},
	}

	p := UGCPolicy()
	p.AllowStyles("color").Globally()
	p.RequireParseableURLs(true)

	// These tests are run concurrently to enable the race detector to pick up
	// potential issues
	wg := sync.WaitGroup{}
	wg.Add(len(tests))
	for ii, tt := range tests {
		go func(ii int, tt test) {
			out := p.Sanitize(tt.in)
			if out != tt.expected {
				t.Errorf(
					"test %d failed;\ninput   : %s\noutput  : %s\nexpected: %s",
					ii,
					tt.in,
					out,
					tt.expected,
				)
			}
			wg.Done()
		}(ii, tt)
	}
	wg.Wait()
}

func TestMatchingHandler(t *testing.T) {
	truthHandler := func(value string) bool {
		return true
	}

	tests := []test{
		{
			in:       `<div style="color: invalidValue"></div>`,
			expected: `<div style="color: invalidValue"></div>`,
		},
	}

	p := UGCPolicy()
	p.AllowStyles("color").MatchingHandler(truthHandler).Globally()
	p.RequireParseableURLs(true)

	// These tests are run concurrently to enable the race detector to pick up
	// potential issues
	wg := sync.WaitGroup{}
	wg.Add(len(tests))
	for ii, tt := range tests {
		go func(ii int, tt test) {
			out := p.Sanitize(tt.in)
			if out != tt.expected {
				t.Errorf(
					"test %d failed;\ninput   : %s\noutput  : %s\nexpected: %s",
					ii,
					tt.in,
					out,
					tt.expected,
				)
			}
			wg.Done()
		}(ii, tt)
	}
	wg.Wait()
}

func TestStyleBlockHandler(t *testing.T) {
	truthHandler := func(value string) bool {
		return true
	}

	tests := []test{
		{
			in:       ``,
			expected: ``,
		},
	}

	p := UGCPolicy()
	p.AllowStyles("color").MatchingHandler(truthHandler).Globally()
	p.RequireParseableURLs(true)

	// These tests are run concurrently to enable the race detector to pick up
	// potential issues
	wg := sync.WaitGroup{}
	wg.Add(len(tests))
	for ii, tt := range tests {
		go func(ii int, tt test) {
			out := p.Sanitize(tt.in)
			if out != tt.expected {
				t.Errorf(
					"test %d failed;\ninput   : %s\noutput  : %s\nexpected: %s",
					ii,
					tt.in,
					out,
					tt.expected,
				)
			}
			wg.Done()
		}(ii, tt)
	}
	wg.Wait()
}
