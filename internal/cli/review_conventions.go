package cli

const deliveryConventionsVersion = "roundfix/delivery-conventions/v1"

type deliveryConvention struct {
	ID     string
	Prompt string
}

func deliveryConventions() []deliveryConvention {
	return []deliveryConvention{
		{"C1", "A Spec's QA Report records the head it audited and is committed after that head, so it never names the commit that records it."},
		{"C2", "The Daemon writes a Task file's status and its `## Result`, `## Recorded paths` and `## Carry-forward provenance` sections after the Task's Verification passes; a Result that calls status Daemon-owned agrees with a `completed` status."},
		{"C3", "The archive commit moves a completed Spec's directory to the archive root and stamps its archive front matter."},
		{"C4", "A planning candidate authors a Spec whose Tasks are all pending and which has no QA Report; that Spec's own delivery implements it and is reviewed then."},
	}
}
