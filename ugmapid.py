#!/usr/bin/env python3

import argparse
import os
import stat


# walk_files performs the actual uid/gid change
def walk_files(root: str, reference, base: int, verbose: bool = True) -> None:
    for basedir, dirs, files in os.walk(root):
        files_map = map(lambda f: os.path.join(basedir, f), files)
        directories_map = map(lambda d: os.path.join(basedir, d), dirs)
        for path in list(files_map) + list(directories_map) + [root]:
            st = os.lstat(path)
            tgt_uid, tgt_gid = st.st_uid, st.st_gid
            if tgt_uid < base or tgt_uid > base + 65536:
                tgt_uid = base + (tgt_uid - reference)
            if tgt_gid < base or tgt_gid > base + 65536:
                tgt_gid = base + (tgt_gid - reference)
            if tgt_uid != st.st_uid or tgt_gid != st.st_gid:
                os.lchown(path, tgt_uid, tgt_gid)
                if stat.S_ISDIR(st.st_mode) or stat.S_ISREG(st.st_mode):
                    os.chmod(path, st.st_mode)
                if verbose:
                    print('{0}: uid {1} -> {2}, gid {3} -> {4}'.format(path, st.st_uid, tgt_uid, st.st_gid, tgt_gid))


# main program
if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('--path', type=str, required=True, help='Root directory to process')
    parser.add_argument('--base', type=int, default=100000, help='Uid/Gid base offset')
    parser.add_argument('--verbose', action='store_true', help='Verbose output')
    args = parser.parse_args()

    # update base
    base_id = os.stat(args.path).st_uid
    if args.verbose:
        print('Calculated base offset: {0}'.format(base_id))
    walk_files(args.path, base_id, args.base, args.verbose)
