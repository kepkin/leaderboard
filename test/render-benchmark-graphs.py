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

fig = plt.figure()
for path in glob.glob(os.path.join(sys.argv[1], "insertion-*.csv")):
	workDir = os.path.dirname(path)
	filename = os.path.basename(path)
	filenameWithoutExt = filename.split(".")[0]
	operation, struct, data = filenameWithoutExt.split('-', 3)
	msft = pd.read_csv(path, index_col=0)

	if len(sys.argv) == 3: 
		msft = msft[sys.argv[2]]

	print(msft)

	fig.clear()
	ax1 = fig.add_subplot()
	ax1.plot(msft.index, msft['tps']) #, grid=True)
	# ax1.xaxis.set_major_formatter(ticks_x)
	ax1.set_ylabel("tps")
	ax1.set_xlabel("total insertions")
	ax1.set_title(operation)
	ax1.grid(True)
	#msft.plot(ylabel="tps", xlabel="total insertions", title=operation, grid=True)

	plt.savefig(os.path.join(workDir, filenameWithoutExt + ".png"))