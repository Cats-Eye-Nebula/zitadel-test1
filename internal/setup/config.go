package setup

import (
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/domain"
)

type IAMSetUp struct {
	Step1  *command.Step1
	Step2  *command.Step2
	Step3  *command.Step3
	Step4  *command.Step4
	Step5  *command.Step5
	Step6  *command.Step6
	Step7  *command.Step7
	Step8  *command.Step8
	Step9  *command.Step9
	Step10 *command.Step10
	Step11 *command.Step11
	Step12 *command.Step12
	Step13 *command.Step13
	Step14 *command.Step14
	Step15 *command.Step15
	Step16 *command.Step16
	Step17 *command.Step17
	Step18 *command.Step18
	Step19 *command.Step19
	Step20 *command.Step20
	Step21 *command.Step21
}

func (setup *IAMSetUp) Steps(currentDone domain.Step) ([]command.Step, error) {
	steps := make([]command.Step, 0)

	for _, step := range []command.Step{
		setup.Step1,
		setup.Step2,
		setup.Step3,
		setup.Step4,
		setup.Step5,
		setup.Step6,
		setup.Step7,
		setup.Step8,
		setup.Step9,
		setup.Step10,
		setup.Step11,
		setup.Step12,
		setup.Step13,
		setup.Step14,
		setup.Step15,
		setup.Step16,
		setup.Step17,
		setup.Step18,
		setup.Step19,
		setup.Step20,
		setup.Step21,
	} {
		if step.Step() <= currentDone {
			continue
		}
		steps = append(steps, step)
	}
	return steps, nil
}
