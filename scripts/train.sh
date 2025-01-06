#!/bin/zsh

export GOMAXPROCS=10

BINDIR=~/work/GitHub/evodevo3/bin
genenv=${BINDIR}/genenv
runsim=${BINDIR}/runsim
gpplot=${BINDIR}/gpplot

ENVS=false
TRAIN=true
TEST=true
ANAL=true

EI1=10 # epochs for training: [0, EI1)
EI2=12 # epochs for testing (production run) [EI1, EI2)

POPSIZE=500
ENVSFILE=Environments.json
if $ENVS; then
    ${genenv} -o ${ENVSFILE} -n 50 -denv 0.5 -seed 13 2> environ.log
fi

EI0=`printf "%2.2d" $[${EI1}-1]`
EEND=`printf "%2.2d" $[${EI2}-1]`

MODELS=(Full NoCue NoDev)
#MODELS=(Hie0 Hie1 Hie2)
#MODELS=(NullCue NullDev NullHie)

if $TRAIN; then
    for model in $MODELS; do
	${runsim} -envs ${ENVSFILE} -model ${model} -popsize ${POPSIZE} \
		  -env_start 0 -env_end ${EI1} -ngen 200 \
		  > data/${model}_train.out
    done 2> train.log
fi

if $TEST; then
    for model in $MODELS; do
	${runsim} -envs ${ENVSFILE} -production \
		  -env_start ${EI1} -env_end ${EI2} \
		  -ngen 200 -setting traj/Setting_${model}.json \
		  -restart traj/${model}_${EI0}_200.traj.gz \
		  > data/${model}_test.out
    done 2> test.log
fi

if $ANAL; then
    for model in $MODELS; do
	for epoch in {${EI1}..${EEND}}; do
	    echo epoch $epoch
	    ${gpplot} -setting traj/Setting_${model}.json \
 		      -envs ${ENVSFILE} -ienv ${epoch} \
		      traj/${model}_${epoch}_*.traj.gz
	done
    done 2> plot.log
fi
