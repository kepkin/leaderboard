import matplotlib.pyplot as plt
import matplotlib.cbook as cbook
import matplotlib.axes as axes
import matplotlib.ticker as ticker

import numpy as np
import pandas as pd
import sys
import glob
import os.path
import humanize as hu

# ticks_x = ticker.FuncFormatter(lambda x, pos: hu.intword(x, format='%.1f'))
# ax = axes.Axes(plt.figure(), ylabel="tps", xlabel="total insertions")

benchmarkReportFolder = sys.argv[1]


fig = plt.figure()
ax = fig.add_subplot()
insertionDF = pd.read_csv("./data/insertion.csv", index_col=0)
ax.hist(insertionDF, bins=600, histtype='step')

ax.legend(insertionDF.columns)

plt.savefig(os.path.join(benchmarkReportFolder, "input-data.png"))
fig.clear()

for path in glob.glob(os.path.join(benchmarkReportFolder, "insertion-*.csv")):
	benchmarkReportFolder = os.path.dirname(path)
	filename = os.path.basename(path)
	filenameWithoutExt = filename.split(".")[0]
	operation, struct, data = filenameWithoutExt.split('-', 3)
	dataFrame = pd.read_csv(path, index_col=0)

	fig.clear()
	ax1 = fig.add_subplot()
	ax1.plot(dataFrame.index, dataFrame['tps']) #, grid=True)
	# ax1.xaxis.set_major_formatter(ticks_x)
	ax1.set_ylabel("tps")
	ax1.set_xlabel("total insertions")
	ax1.set_title("{} - {} - {}".format(operation, struct, data))
	ax1.grid(True)

	plt.savefig(os.path.join(benchmarkReportFolder, filenameWithoutExt + ".png"))