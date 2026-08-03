BEGIN {
	FS = "\t"
}

NF == 2 {
	count[$1]++
	total++
}

NF != 2 {
	malformed++
}

END {
	production = count["rev-parse"] + count["ls-tree"] + count["cat-file"] + count["rev-list"] + count["status"] + count["for-each-ref"]
	fixture = count["init"] + count["config"] + count["add"] + count["commit"]
	ambiguous = total - production - fixture

	printf "%d production-read-shaped\n", production
	printf "%d fixture-setup-shaped\n", fixture
	printf "%d ambiguous-or-other\n", ambiguous
	printf "%d malformed-records\n", malformed
}
