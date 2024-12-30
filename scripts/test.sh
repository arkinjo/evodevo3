#!/bin/zsh

export GOMAXPROCS=8

BINDIR=~/work/GitHub/evodevo3/bin
runsim=${BINDIR}/runsim
gpplot=${BINDIR}/gpplot

POPSIZE=500
ENVSFILE=Environments.json
EI1=30
EI0=`printf "%2.2d" $[${EI1}-1]`
EI2=40
EEND=`printf "%2.2d" $[${EI2}-1]`
MODELS=(FULL)
#MODELS=(Full NoCue NoDev)
#MODELS=(Hie0 Hie1 Hie2)
#MODELS=(NullCue NullDev NullHie)
for model in $MODELS; do
    ${runsim} -envs ${ENVSFILE} -production -env_start ${EI1} -env_end ${EI2} \
	      -ngen 200 -setting traj/Setting_${model}.json \
	      -restart traj/${model}_${EI0}_200.traj.gz \
	       > data/${model}_test.out
done 2> test.log
