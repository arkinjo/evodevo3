#!/bin/zsh

export GOMAXPROCS=8

BINDIR=~/work/GitHub/evodevo3/bin
runsim=${BINDIR}/runsim
gpplot=${BINDIR}/gpplot

POPSIZE=500
ENVSFILE=Environments.json
EI1=30
EI0=`printf "%2.2d" $[${EI1}-1]`
EI2=31
EEND=`printf "%2.2d" $[${EI2}-1]`
MODELS=(FULL)
#MODELS=(Full NoCue NoDev)
#MODELS=(Hie0 Hie1 Hie2)
#MODELS=(NullCue NullDev NullHie)
for model in $MODELS; do
    for epoch in {${EI1}..${EEND}}; do
	echo epoch $epoch
	${gpplot} -setting traj/Setting_${model}.json \
 		  -envs ${ENVSFILE} -ienv ${epoch} traj/${model}_${epoch}_*.traj.gz
    done
done 2> plot.log
