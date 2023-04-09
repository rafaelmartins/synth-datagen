package codegen

import (
	"fmt"
	"io"
)

type Header struct {
	include includeList
	macro   macroList
	data    dataList
}

func NewHeader() *Header {
	return &Header{
		include: includeList{},
		macro:   macroList{},
		data:    dataList{},
	}
}

func (h *Header) AddInclude(path string, system bool) {
	h.include.add(path, system)
}

func (h *Header) AddMacro(identifier string, value interface{}, hex bool, raw bool) {
	h.macro.add(identifier, value, hex, raw)
}

func (h *Header) AddData(identifier string, value interface{}, attributes []string) {
	h.data.add(identifier, value, attributes)
}

func (h *Header) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, `// Code generated by "synth-datagen"; DO NOT EDIT.

// SPDX-FileCopyrightText: 2022-present Rafael G. Martins <rafael@rafaelmartins.eng.br>
// SPDX-License-Identifier: BSD-3-Clause

#pragma once
`); err != nil {
		return err
	}

	if err := h.include.write(w); err != nil {
		return err
	}

	if err := h.macro.write(w); err != nil {
		return err
	}

	if err := h.data.write(w); err != nil {
		return err
	}

	return nil
}
