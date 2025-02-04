// Copyright (C) 2022 Red Hat, Inc.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program; if not, write to the Free Software Foundation, Inc.,
// 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.

package tolerations

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
)

//nolint:funlen
func TestIsTolerationModified(t *testing.T) {
	getInt64Pointer := func(val int64) *int64 {
		return &val
	}

	testCases := []struct {
		testToleration v1.Toleration
		expectedOutput bool
		qosClass       v1.PodQOSClass
	}{
		{ // Test Case #1 - default not-ready toleration
			testToleration: v1.Toleration{
				Key:               "node.kubernetes.io/not-ready",
				Operator:          v1.TolerationOpExists,
				Effect:            v1.TaintEffectNoExecute,
				TolerationSeconds: getInt64Pointer(300),
			},
			expectedOutput: false,
			qosClass:       v1.PodQOSGuaranteed,
		},
		{ // Test Case #2 - default unreachable toleration
			testToleration: v1.Toleration{
				Key:               "node.kubernetes.io/unreachable",
				Operator:          v1.TolerationOpExists,
				Effect:            v1.TaintEffectNoExecute,
				TolerationSeconds: getInt64Pointer(300),
			},
			expectedOutput: false,
			qosClass:       v1.PodQOSGuaranteed,
		},
		{ // Test Case #3 - modified unreachable toleration
			testToleration: v1.Toleration{
				Key:               "node.kubernetes.io/unreachable",
				Operator:          v1.TolerationOpExists,
				Effect:            v1.TaintEffectNoExecute,
				TolerationSeconds: getInt64Pointer(350), // modified from 300
			},
			expectedOutput: true,
			qosClass:       v1.PodQOSGuaranteed,
		},
		{ // Test Case #4 - modified unreachable toleration
			testToleration: v1.Toleration{
				Key:               "node.kubernetes.io/unreachable",
				Operator:          v1.TolerationOpEqual, // modified from exists
				Effect:            v1.TaintEffectNoExecute,
				TolerationSeconds: getInt64Pointer(300),
			},
			expectedOutput: true,
			qosClass:       v1.PodQOSGuaranteed,
		},
		{ // Test Case #5 - missing effect
			testToleration: v1.Toleration{
				Key:      "node.kubernetes.io/unreachable",
				Operator: v1.TolerationOpExists,
				// Effect:            v1.TaintEffectNoExecute,
				TolerationSeconds: getInt64Pointer(300),
			},
			expectedOutput: false,
			qosClass:       v1.PodQOSGuaranteed,
		},
		{ // Test Case #6 - example from QE and DCI - this should pass only if qosClass is
			// different than BestEffort, which is the case
			testToleration: v1.Toleration{
				Key:      "node.kubernetes.io/memory-pressure",
				Operator: v1.TolerationOpExists,
				Effect:   v1.TaintEffectNoSchedule,
			},
			expectedOutput: false,
			qosClass:       v1.PodQOSGuaranteed,
		},
		{ // Test Case #7 - example from QE and DCI - however, if qosClass is BestEffort, it
			// must be considered as a modified toleration
			testToleration: v1.Toleration{
				Key:      "node.kubernetes.io/memory-pressure",
				Operator: v1.TolerationOpExists,
				Effect:   v1.TaintEffectNoSchedule,
			},
			expectedOutput: true,
			qosClass:       v1.PodQOSBestEffort,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expectedOutput, IsTolerationModified(tc.testToleration, tc.qosClass))
	}
}
