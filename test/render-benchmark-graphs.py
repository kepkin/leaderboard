#!/usr/bin/env python3

from fileinput import filename
import matplotlib.pyplot as plt
import matplotlib.cbook as cbook
import matplotlib.axes as axes
import matplotlib.ticker as ticker

import argparse
import numpy as np
import pandas as pd
import sys
import glob
import os.path
import humanize as hu

# ticks_x = ticker.FuncFormatter(lambda x, pos: hu.intword(x, format='%.1f'))
# ax = axes.Axes(plt.figure(), ylabel="tps", xlabel="total insertions")

fig = plt.figure()


def plotInputInsertions(targetFolder):
	ax = fig.add_subplot()
	insertionDF = pd.read_csv("./data/insertion.csv", index_col=0)
	ax.hist(insertionDF, bins=600, histtype='step')

	ax.legend(insertionDF.columns)

	plt.savefig(os.path.join(targetFolder, "input-data.png"))


def GetFilenameWithoutExt(path):
	filename = os.path.basename(path)
	return ".".join(filename.split(".")[:-1])

class InsertionPlotBuilder:
	def __init__(self, path):
		self.path = path

	def GetInsertionMetaInfo(self):
		filenameWithoutExt = GetFilenameWithoutExt(self.path)
		operation, struct, data = filenameWithoutExt.split('-', 3)
		return operation, struct, data

	def BuildPlot(self, ax, title=None, label=None):
		if title is None:
			operation, struct, data = self.GetInsertionMetaInfo()
			title = "{} - {} - {}".format(operation, struct, data)

		dataFrame = pd.read_csv(self.path, index_col=0)

		ax.plot(dataFrame.index, dataFrame['tps'], label=label)
		ax.set_ylabel("tps")
		ax.set_xlabel("total insertions")
		ax.set_title(title)
		ax.grid(True)
		ax.legend()


def plotInsertionPlots(targetFolder, comparingFolder=None):
	def buildHelper(path, label=None):
		InsPlotBuilder = InsertionPlotBuilder(path)
		InsPlotBuilder.BuildPlot(fig.add_subplot(), label=label)


	for path in glob.glob(os.path.join(targetFolder, "insertion-*.csv")):
		fig.clear()
		filename = os.path.basename(path)

		path2 = None
		if comparingFolder:
			path2 = os.path.join(comparingFolder, filename)
		
		if path2 and os.path.exists(path2):
			buildHelper(path, "now")
			buildHelper(path2, "before")
		else:
			buildHelper(path)

		filenameWithoutExt = GetFilenameWithoutExt(path)
		plt.savefig(os.path.join(targetFolder, filenameWithoutExt + ".png"))

def compareTwoIsertionCSVs(pathA, labelA, pathB, labelB, targetFile):
	fig.clear()
	InsPlotBuilder = InsertionPlotBuilder(pathA)
	InsPlotBuilder.BuildPlot(fig.add_subplot(), label=labelA)
	InsPlotBuilder = InsertionPlotBuilder(pathB)
	InsPlotBuilder.BuildPlot(fig.add_subplot(), label=labelB)
	plt.savefig(targetFile + ".png")

	
def plotBenchmarkFolder(args):

	plotInputInsertions(args.path)
	plotInsertionPlots(args.path, args.path2)

def cmpInsertionCSV(args):
	InsPlotBuilderA = InsertionPlotBuilder(args.patha)
	InsPlotBuilderB = InsertionPlotBuilder(args.pathb)

	labelA = "{m[1]} - {m[2]}".format(m=InsPlotBuilderA.GetInsertionMetaInfo())
	labelB = "{m[1]} - {m[2]}".format(m=InsPlotBuilderB.GetInsertionMetaInfo())

	compareTwoIsertionCSVs(args.patha, labelA, args.pathb, labelB, args.target)

parser = argparse.ArgumentParser(prog='PROG')

subparsers = parser.add_subparsers(help='sub-command help')
# create the parser for the "a" command
parser_a = subparsers.add_parser('plotBenchmarkFolder', help='a help')
parser_a.set_defaults(func=plotBenchmarkFolder)

parser_a.add_argument('path', help='path to folder')
parser_a.add_argument('path2', help='path to folder for building comparions graphs')
# create the parser for the "b" command
parser_b = subparsers.add_parser('cmpInsertionCSV', help='compare two specific insertion CSVs')
parser_b.set_defaults(func=cmpInsertionCSV)

parser_b.add_argument('patha', help='path to a csv')
parser_b.add_argument('pathb', help='path to b csv')
parser_b.add_argument('target', help='path to target image without EXT')

# parse some argument lists
args = parser.parse_args()
args.func(args)
