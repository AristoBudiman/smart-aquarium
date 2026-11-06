// Kode dibuat oleh ml/train_export.py. Jangan diedit.
package model

func ScoreKekeruhan(input []float64) []float64 {
    var var0 []float64
    if input[0] <= 24.975000381469727 {
        var0 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.119998931884766 {
            var0 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.37999725341797 {
                var0 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var0 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var1 []float64
    if input[0] <= 24.994999885559082 {
        var1 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.09499931335449 {
            var1 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.68000030517578 {
                var1 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var1 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var2 []float64
    if input[0] <= 24.984999656677246 {
        var2 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.09499931335449 {
            var2 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.56999969482422 {
                var2 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var2 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var3 []float64
    if input[0] <= 24.984999656677246 {
        var3 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.119998931884766 {
            var3 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.68000030517578 {
                var3 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var3 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var4 []float64
    if input[0] <= 24.985000610351562 {
        var4 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.15999984741211 {
            var4 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.68000030517578 {
                var4 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var4 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var5 []float64
    if input[0] <= 24.994999885559082 {
        var5 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.18499946594238 {
            var5 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.95999908447266 {
                var5 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var5 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var6 []float64
    if input[0] <= 24.984999656677246 {
        var6 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.119998931884766 {
            var6 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.68000030517578 {
                var6 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var6 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var7 []float64
    if input[0] <= 24.984999656677246 {
        var7 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.18499946594238 {
            var7 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.37999725341797 {
                var7 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var7 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var8 []float64
    if input[0] <= 24.984999656677246 {
        var8 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.18499946594238 {
            var8 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.95999908447266 {
                var8 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var8 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var9 []float64
    if input[0] <= 24.94499969482422 {
        var9 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.10000038146973 {
            var9 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 151.9800033569336 {
                var9 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var9 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var10 []float64
    if input[0] <= 24.99000072479248 {
        var10 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.119998931884766 {
            var10 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.68000030517578 {
                var10 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var10 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var11 []float64
    if input[0] <= 24.984999656677246 {
        var11 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.18499946594238 {
            var11 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.68000030517578 {
                var11 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var11 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var12 []float64
    if input[0] <= 24.994999885559082 {
        var12 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.119998931884766 {
            var12 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.95999908447266 {
                var12 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var12 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var13 []float64
    if input[0] <= 24.994999885559082 {
        var13 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.204999923706055 {
            var13 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.33999633789062 {
                var13 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var13 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var14 []float64
    if input[0] <= 24.975000381469727 {
        var14 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.119998931884766 {
            var14 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.68000030517578 {
                var14 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var14 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var15 []float64
    if input[0] <= 24.984999656677246 {
        var15 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.119998931884766 {
            var15 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.68000030517578 {
                var15 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var15 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var16 []float64
    if input[0] <= 24.975000381469727 {
        var16 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.09499931335449 {
            var16 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 151.68000030517578 {
                var16 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var16 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var17 []float64
    if input[0] <= 24.994999885559082 {
        var17 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.18499946594238 {
            var17 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.95999908447266 {
                var17 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var17 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var18 []float64
    if input[0] <= 24.984999656677246 {
        var18 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.204999923706055 {
            var18 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.68000030517578 {
                var18 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var18 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var19 []float64
    if input[0] <= 24.994999885559082 {
        var19 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.204999923706055 {
            var19 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.65999603271484 {
                var19 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var19 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var20 []float64
    if input[0] <= 24.984999656677246 {
        var20 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.18499946594238 {
            var20 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.68000030517578 {
                var20 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var20 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var21 []float64
    if input[0] <= 24.984999656677246 {
        var21 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.18499946594238 {
            var21 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.68000030517578 {
                var21 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var21 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var22 []float64
    if input[0] <= 24.975000381469727 {
        var22 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.18499946594238 {
            var22 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.95999908447266 {
                var22 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var22 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var23 []float64
    if input[0] <= 24.994999885559082 {
        var23 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.119998931884766 {
            var23 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 151.9800033569336 {
                var23 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var23 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var24 []float64
    if input[0] <= 24.994999885559082 {
        var24 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.15999984741211 {
            var24 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.68000030517578 {
                var24 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var24 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var25 []float64
    if input[0] <= 24.994999885559082 {
        var25 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.15999984741211 {
            var25 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.37999725341797 {
                var25 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var25 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var26 []float64
    if input[0] <= 24.984999656677246 {
        var26 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.2450008392334 {
            var26 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.95999908447266 {
                var26 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var26 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var27 []float64
    if input[0] <= 24.984999656677246 {
        var27 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.119998931884766 {
            var27 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.33999633789062 {
                var27 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var27 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var28 []float64
    if input[0] <= 24.984999656677246 {
        var28 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.09499931335449 {
            var28 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.6199951171875 {
                var28 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var28 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var29 []float64
    if input[0] <= 24.984999656677246 {
        var29 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.119998931884766 {
            var29 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.65999603271484 {
                var29 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var29 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var30 []float64
    if input[0] <= 24.994999885559082 {
        var30 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.119998931884766 {
            var30 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.68000030517578 {
                var30 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var30 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var31 []float64
    if input[0] <= 24.994999885559082 {
        var31 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.27000045776367 {
            var31 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.31500244140625 {
                var31 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var31 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var32 []float64
    if input[0] <= 24.984999656677246 {
        var32 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.18499946594238 {
            var32 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.68000030517578 {
                var32 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var32 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var33 []float64
    if input[0] <= 24.994999885559082 {
        var33 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.09499931335449 {
            var33 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.95999908447266 {
                var33 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var33 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var34 []float64
    if input[0] <= 24.994999885559082 {
        var34 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.18000030517578 {
            var34 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 153.00499725341797 {
                var34 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var34 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var35 []float64
    if input[0] <= 24.994999885559082 {
        var35 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.18000030517578 {
            var35 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.68000030517578 {
                var35 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var35 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var36 []float64
    if input[0] <= 24.960000038146973 {
        var36 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.18499946594238 {
            var36 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 151.9800033569336 {
                var36 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var36 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var37 []float64
    if input[0] <= 24.994999885559082 {
        var37 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.09499931335449 {
            var37 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.65999603271484 {
                var37 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var37 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var38 []float64
    if input[0] <= 24.994999885559082 {
        var38 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.119998931884766 {
            var38 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.59500122070312 {
                var38 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var38 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var39 []float64
    if input[0] <= 24.994999885559082 {
        var39 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.18000030517578 {
            var39 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.68000030517578 {
                var39 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var39 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var40 []float64
    if input[0] <= 25.019999504089355 {
        var40 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.119998931884766 {
            var40 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 152.10499572753906 {
                var40 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var40 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var41 []float64
    if input[0] <= 24.994999885559082 {
        var41 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.09499931335449 {
            var41 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.31500244140625 {
                var41 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var41 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var42 []float64
    if input[0] <= 24.984999656677246 {
        var42 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.09499931335449 {
            var42 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.68000030517578 {
                var42 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var42 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var43 []float64
    if input[0] <= 24.984999656677246 {
        var43 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.09499931335449 {
            var43 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.68000030517578 {
                var43 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var43 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var44 []float64
    if input[0] <= 24.994999885559082 {
        var44 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.18499946594238 {
            var44 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 151.68000030517578 {
                var44 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var44 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var45 []float64
    if input[0] <= 24.97000026702881 {
        var45 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.119998931884766 {
            var45 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.68000030517578 {
                var45 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var45 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var46 []float64
    if input[0] <= 24.984999656677246 {
        var46 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.18000030517578 {
            var46 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 151.68000030517578 {
                var46 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var46 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var47 []float64
    if input[0] <= 24.984999656677246 {
        var47 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.119998931884766 {
            var47 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.68000030517578 {
                var47 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var47 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var48 []float64
    if input[0] <= 24.984999656677246 {
        var48 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.119998931884766 {
            var48 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.68000030517578 {
                var48 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var48 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var49 []float64
    if input[0] <= 24.994999885559082 {
        var49 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 50.18499946594238 {
            var49 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 150.65999603271484 {
                var49 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var49 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    return mulVectorNumberKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(addVectorsKekeruhan(var0, var1), var2), var3), var4), var5), var6), var7), var8), var9), var10), var11), var12), var13), var14), var15), var16), var17), var18), var19), var20), var21), var22), var23), var24), var25), var26), var27), var28), var29), var30), var31), var32), var33), var34), var35), var36), var37), var38), var39), var40), var41), var42), var43), var44), var45), var46), var47), var48), var49), 0.02)
}
func addVectorsKekeruhan(v1, v2 []float64) []float64 {
    result := make([]float64, len(v1))
    for i := 0; i < len(v1); i++ {
        result[i] = v1[i] + v2[i]
    }
    return result
}
func mulVectorNumberKekeruhan(v1 []float64, num float64) []float64 {
    result := make([]float64, len(v1))
    for i := 0; i < len(v1); i++ {
        result[i] = v1[i] * num
    }
    return result
}
